package download

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/collection"
	"github.com/fekoneko/piximan/internal/collection/work"
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
	paths := utils.FromPtr(options.Paths, []string{""})

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
		d.Schedule(*options.Ids, kind, size, onlyMeta, paths)

	} else if options.Bookmarks != nil && *options.Bookmarks == "my" {
		d.ScheduleMyBookmarks(
			kind, options.Tags, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, untilSkip, paths,
		)

	} else if options.Bookmarks != nil {
		userId, err := strconv.ParseUint(*options.Bookmarks, 10, 64)
		logger.MaybeFatal(err, "cannot parse user id %v", *options.Bookmarks)

		d.ScheduleBookmarks(
			userId, kind, options.Tags, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, untilSkip, paths,
		)

	} else if options.InferIds != nil {
		q := make(queue.Queue, 0)
		mutex := &sync.Mutex{}
		waitGroup := &sync.WaitGroup{}
		seen := make(map[string]bool, len(*options.InferIds))

		for _, rawInferId := range *options.InferIds {
			inferId := filepath.Clean(rawInferId)
			if seen[inferId] {
				continue
			}
			seen[inferId] = true
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				var items []queue.Item

				if fsext.IsInferIdPattern(inferId) {
					idPathsMap, errs := fsext.InferIdsFromPattern(inferId)
					logger.MaybeErrors(errs, "error while inferring work id from pattern %v", inferId)
					if len(*idPathsMap) == 0 {
						logger.Fatal("no ids could be inferred from pattern %v", inferId)
					}

					items = make([]queue.Item, 0, len(*idPathsMap))
					for id, paths := range *idPathsMap {
						items = append(items, queue.Item{
							Id:       id,
							Kind:     kind,
							Size:     size,
							OnlyMeta: onlyMeta,
							Paths:    utils.If(options.Paths != nil, paths, paths),
						})
					}

				} else {
					c := collection.New(inferId, logger.DefaultLogger)
					// TODO: collection.ReadQueue() that will only care about id and kind in metadata.yaml and ignore assets
					c.Read()
					artworksIdPathsMap := utils.ToPtr(make(map[uint64][]string))
					novelsIdPathsMap := utils.ToPtr(make(map[uint64][]string))
					for w := c.WaitNext(); w != nil; w = c.WaitNext() {
						if w.Id != nil && w.Kind != nil {
							// TODO: work.Kind.IsArtwork() / work.Kind.IsNovel()
							switch *w.Kind {
							case work.KindIllust, work.KindManga, work.KindUgoira:
								(*artworksIdPathsMap)[*w.Id] = append((*artworksIdPathsMap)[*w.Id], w.Path)
							case work.KindNovel:
								(*novelsIdPathsMap)[*w.Id] = append((*novelsIdPathsMap)[*w.Id], w.Path)
							}
						}
					}
					if options.Kind != nil && kind == queue.ItemKindArtwork && len(*artworksIdPathsMap) == 0 {
						logger.Fatal("no works with id and artwork kind found in directory %v", inferId)
					} else if options.Kind != nil && kind == queue.ItemKindNovel && len(*novelsIdPathsMap) == 0 {
						logger.Fatal("no works with id and novel kind found in directory %v", inferId)
					} else if len(*artworksIdPathsMap) == 0 && len(*novelsIdPathsMap) == 0 {
						logger.Fatal("no works with id and kind found in directory %v", inferId)
					}

					items = make([]queue.Item, 0, len(*artworksIdPathsMap))
					if options.Kind == nil || kind == queue.ItemKindArtwork {
						for id, inferredPaths := range *artworksIdPathsMap {
							items = append(items, queue.Item{
								Id:       id,
								Kind:     queue.ItemKindArtwork,
								Size:     size,
								OnlyMeta: onlyMeta,
								Paths:    utils.If(options.Paths != nil, paths, inferredPaths),
							})
						}
					}
					if options.Kind == nil || kind == queue.ItemKindNovel {
						for id, inferredPaths := range *novelsIdPathsMap {
							items = append(items, queue.Item{
								Id:       id,
								Kind:     queue.ItemKindNovel,
								Size:     size,
								OnlyMeta: onlyMeta,
								Paths:    utils.If(options.Paths != nil, paths, inferredPaths),
							})
						}
					}
				}

				logger.Info(
					"inferred %v work%v from directory %v",
					len(items), utils.Plural(len(items)), inferId,
				)
				mutex.Lock()
				q = append(q, items...)
				mutex.Unlock()
			}()
		}

		waitGroup.Wait()
		d.ScheduleQueue(&q)

	} else if options.Lists != nil {
		seen := make(map[string]bool, len(*options.Lists))

		for _, rawListPath := range *options.Lists {
			listPath := filepath.Clean(rawListPath)
			if seen[listPath] {
				continue
			}
			seen[listPath] = true

			q, err := fsext.ReadList(listPath, kind, size, onlyMeta, paths)
			logger.MaybeFatal(err, "cannot read download list from %v", listPath)
			if len(*q) == 0 {
				logger.Fatal("no works found in list %v", listPath)
			}
			d.ScheduleQueue(q)
		}
	}

	if options.Rules != nil {
		for _, rulesPath := range *options.Rules {
			rules, err := fsext.ReadRules(rulesPath)
			logger.MaybeFatal(err, "cannot read download rules from %v", rulesPath)
			d.AddRules(*rules)
		}
	}

	if options.Skips != nil && len(*options.Skips) != 0 {
		list := skiplist.New()
		mutex := &sync.Mutex{}
		waitGroup := &sync.WaitGroup{}
		seen := make(map[string]bool, len(*options.Skips))

		for _, rawSkipPath := range *options.Skips {
			skipPath := filepath.Clean(rawSkipPath)
			if seen[skipPath] {
				continue
			}
			seen[skipPath] = true
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				numWorks := 0
				if fsext.IsInferIdPattern(skipPath) {
					idPathMap, errs := fsext.InferIdsFromPattern(skipPath)
					logger.MaybeErrors(errs, "error while inferring work id from pattern %v", skipPath)
					numWorks = len(*idPathMap)
					if numWorks == 0 {
						logger.Fatal("no ids could be inferred from pattern %v", skipPath)
					}
					mutex.Lock()
					for id := range *idPathMap {
						list.Add(id, kind)
					}
					mutex.Unlock()

				} else {
					c := collection.New(skipPath, logger.DefaultLogger)
					// TODO: collection.ReadQueue() that will only care about id and kind in metadata.yaml and ignore assets
					c.Read()
					for w := c.WaitNext(); w != nil; w = c.WaitNext() {
						if w.Id != nil && w.Kind != nil {
							mutex.Lock()
							list.AddWork(w.Work)
							mutex.Unlock()
							numWorks++
						}
					}
					if numWorks == 0 {
						logger.Fatal("no works with id and kind found in directory %v", skipPath)
					}
				}

				logger.Info(
					"%v work%v found in the directory %v will be skipped",
					numWorks, utils.Plural(numWorks), skipPath,
				)
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
