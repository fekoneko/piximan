package download

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/downloader"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/termext"
	"github.com/fekoneko/piximan/internal/utils"
)

func download(options *options) {
	size := utils.FromPtrTransform(options.Size, image.SizeFromUint, image.SizeDefault)
	kind := utils.FromPtrTransform(options.Kind, queue.ItemKindFromString, queue.ItemKindDefault)
	private := utils.FromPtr(options.Private, false)
	onlyMeta := utils.FromPtr(options.OnlyMeta, false)
	lowMeta := utils.FromPtr(options.LowMeta, false)
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
			kind, options.Tag, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, paths,
		)
		fmt.Println()

	} else if options.Bookmarks != nil {
		userId, err := strconv.ParseUint(*options.Bookmarks, 10, 64)
		logger.MaybeFatal(err, "cannot parse user id %v", *options.Bookmarks)

		paths := []string{path}
		d.ScheduleBookmarks(
			userId, kind, options.Tag, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, paths,
		)
		fmt.Println()

	} else if options.InferId != nil {
		idPathMap, errs := fsext.InferIdsFromWorkPath(*options.InferId)
		logger.MaybeErrors(errs, "error while inferring work id from pattern %v", *options.InferId)
		if len(*idPathMap) == 0 {
			logger.Warning("no ids could be inferred from pattern %v", *options.InferId)
			return
		}

		if options.Path == nil {
			q := queue.FromMap(idPathMap, kind, size, onlyMeta)
			d.ScheduleQueue(q)
		} else {
			paths := []string{path}
			q := queue.FromMapWithPaths(idPathMap, kind, size, onlyMeta, paths)
			d.ScheduleQueue(q)
		}

	} else if options.List != nil {
		paths := []string{path}
		q, err := fsext.ReadList(*options.List, kind, size, onlyMeta, paths)
		logger.MaybeFatal(err, "cannot read the list from %v", *options.List)
		if len(*q) == 0 {
			logger.Warning("no works found in the list %v", *options.List)
			return
		}

		d.ScheduleQueue(q)
	}

	fmt.Println(d)

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
