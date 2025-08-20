package download

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/collection"
	"github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/downloader"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/downloader/skiplist"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/termext"
	"github.com/fekoneko/piximan/internal/utils"
)

func download(options *options) {
	size := utils.FromPtrTransform(options.Size, imageext.SizeFromUint, imageext.SizeDefault)
	kind := utils.FromPtrTransform(options.Kind, queue.ItemKindFromString, queue.ItemKindDefault)
	private := utils.FromPtr(options.Private, false)
	onlyMeta := utils.FromPtr(options.OnlyMeta, false)
	lowMeta := utils.FromPtr(options.LowMeta, false)
	untilSkip := utils.FromPtr(options.UntilSkip, false)
	path := utils.FromPtr(options.Path, "")

	storage, sessionId := configAndSession(options.Password)
	c := client.New(
		sessionId, logger.DefaultLogger,
		storage.PximgMaxPending, storage.PximgDelay,
		storage.DefaultMaxPending, storage.DefaultDelay,
	)
	d := downloader.New(c, logger.DefaultLogger)

	termext.DisableInputEcho()
	defer termext.RestoreInputEcho()

	if options.Ids != nil {
		paths := []string{path}
		d.Schedule(*options.Ids, kind, size, onlyMeta, paths)

	} else if options.Bookmarks != nil && *options.Bookmarks == "my" {
		paths := []string{path}
		d.ScheduleMyBookmarks(
			kind, options.Tags, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, untilSkip, paths,
		)

	} else if options.Bookmarks != nil {
		userId, err := strconv.ParseUint(*options.Bookmarks, 10, 64)
		logger.MaybeFatal(err, "cannot parse user id %v", *options.Bookmarks)

		paths := []string{path}
		d.ScheduleBookmarks(
			userId, kind, options.Tags, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, untilSkip, paths,
		)

	} else if options.InferId != nil {
		q := make(queue.Queue, 0)
		mutex := &sync.Mutex{}
		waitGroup := &sync.WaitGroup{}
		seen := make(map[string]bool, len(*options.InferId))

		for _, rawInferId := range *options.InferId {
			inferId := filepath.Clean(rawInferId)
			if seen[inferId] {
				continue
			}
			seen[inferId] = true
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				withProvidedPaths := options.Path != nil
				providedPaths := []string{path}
				var idPathsMap *map[uint64][]string

				if fsext.IsInferIdPattern(inferId) {
					var errs []error
					idPathsMap, errs = fsext.InferIdsFromPattern(inferId)
					logger.MaybeErrors(errs, "error while inferring work id from pattern %v", inferId)
					if len(*idPathsMap) == 0 {
						logger.Fatal("no ids could be inferred from pattern %v", inferId)
					}

				} else {
					c := collection.New(inferId, logger.DefaultLogger)
					c.Read()
					idPathsMap = utils.ToPtr(make(map[uint64][]string))
					for w := c.WaitNext(); w != nil; w = c.WaitNext() {
						if w.Work.Id != nil {
							(*idPathsMap)[*w.Id] = append((*idPathsMap)[*w.Work.Id], w.Path)
						}
					}
					if len(*idPathsMap) == 0 {
						logger.Fatal("no works found in the directory %v", inferId)
					}
				}

				items := make([]queue.Item, 0, len(*idPathsMap))
				for id, paths := range *idPathsMap {
					items = append(items, queue.Item{
						Id:       id,
						Kind:     kind,
						Size:     size,
						OnlyMeta: onlyMeta,
						Paths:    utils.If(withProvidedPaths, providedPaths, paths),
					})
				}
				mutex.Lock()
				q = append(q, items...)
				mutex.Unlock()
			}()
		}

		waitGroup.Wait()
		d.ScheduleQueue(&q)

	} else if options.List != nil {
		paths := []string{path}
		q, err := fsext.ReadList(*options.List, kind, size, onlyMeta, paths)
		logger.MaybeFatal(err, "cannot read download list from %v", *options.List)
		if len(*q) == 0 {
			logger.Warning("no works found in the list %v", *options.List)
			return
		}
		d.ScheduleQueue(q)
	}

	if options.Rules != nil {
		rules, err := fsext.ReadRules(*options.Rules)
		logger.MaybeFatal(err, "cannot read download rules from %v", *options.Rules)
		d.SetRules(rules)
	}

	if options.Skip != nil && len(*options.Skip) != 0 {
		list := skiplist.New()
		mutex := &sync.Mutex{}
		waitGroup := &sync.WaitGroup{}
		seen := make(map[string]bool, len(*options.Skip))

		for _, rawSkipPath := range *options.Skip {
			skipPath := filepath.Clean(rawSkipPath)
			if seen[skipPath] {
				continue
			}
			seen[skipPath] = true
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				if fsext.IsInferIdPattern(skipPath) {
					idPathMap, errs := fsext.InferIdsFromPattern(skipPath)
					logger.MaybeErrors(errs, "error while inferring work id from pattern %v", skipPath)
					if len(*idPathMap) == 0 {
						logger.Fatal("no ids could be inferred from pattern %v", skipPath)
					}
					mutex.Lock()
					for id := range *idPathMap {
						list.Add(id, kind)
					}
					mutex.Unlock()

				} else {
					c := collection.New(skipPath, logger.DefaultLogger)
					c.Read()
					numWorks := 0
					for w := c.WaitNext(); w != nil; w = c.WaitNext() {
						mutex.Lock()
						list.AddWork(w.Work)
						mutex.Unlock()
						numWorks++
					}
					if numWorks == 0 {
						logger.Fatal("no works found in the directory %v", skipPath)
					}
				}
			}()
		}

		waitGroup.Wait()
		d.SetSkipList(list)
	}

	logger.Info("created downloader:\n%v", d.String())

	logger.EnableProgress()
	defer logger.DisableProgress()

	logger.Info("download started")
	d.Run() // TODO: confirmation by user (or -y flag)
	d.WaitDone()
	logger.Info("download finished")
	logger.Stats()
}

func configAndSession(password *string) (storage *config.Storage, sessionId *string) {
	storage, err := config.New(password)
	if err != nil && password != nil {
		logger.Fatal("cannot open config storage: %v", err)
		panic("unreachable")
	} else if err != nil {
		logger.Warning("cannot open config storage: %v", err)
		promptDefaultConfig()
		return storage, nil
	}

	if sessionId, err := storage.SessionId(); err != nil && password != nil {
		logger.Fatal("cannot read session id: %v", err)
		panic("unreachable")
	} else if err != nil {
		if newStorage, sessionId := promptPassword(); newStorage != nil {
			return newStorage, sessionId
		} else {
			return storage, sessionId
		}
	} else if sessionId == nil && password != nil {
		logger.Fatal("no session id were configured, but password was provided")
		panic("unreachable")
	} else if sessionId == nil {
		logger.Info("no session id were configured, using only anonymous requests")
		return storage, nil
	} else {
		return storage, sessionId
	}
}

func promptPassword() (storage *config.Storage, sessionId *string) {
	for tries := 0; ; tries++ {
		password, err := passwordPrompt.Run()
		if err != nil {
			logger.Warning("failed to read password: %v", err)
			promptNoAuthorization()
			return nil, nil
		}

		storage, err := config.New(&password)
		if err != nil {
			logger.Warning("cannot open config storage: %v", err)
			promptNoAuthorization()
			return nil, nil
		}

		if sessionId, err := storage.SessionId(); err == nil && sessionId != nil {
			return storage, sessionId
		} else if err == nil {
			logger.Info("no session id were configured, using only anonymous requests")
			return storage, nil
		} else if tries == 2 {
			logger.Warning("cannot read session id: %v", err)
			promptNoAuthorization()
			return storage, nil
		}
	}
}

func promptDefaultConfig() {
	_, option, err := deafultConfigPrompt.Run()
	logger.MaybeFatal(err, "failed to read the choice")
	if err != nil || option != YesOption {
		os.Exit(1)
	}
}

func promptNoAuthorization() {
	_, option, err := noAuthorizationPrompt.Run()
	logger.MaybeFatal(err, "failed to read the choice")
	if err != nil || option != YesOption {
		os.Exit(1)
	}
}
