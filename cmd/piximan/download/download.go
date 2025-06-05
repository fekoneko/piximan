package download

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/downloader"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/logext"
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

	config, sessionId := configSession(options.Password)
	d := downloader.New(
		sessionId,
		config.PximgMaxPending, config.PximgDelay,
		config.DefaultMaxPending, config.DefaultDelay,
	)

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
		logext.MaybeFatal(err, "cannot parse user id %v", *options.Bookmarks)

		paths := []string{path}
		d.ScheduleBookmarks(
			userId, kind, options.Tag, options.FromOffset, options.ToOffset, private,
			size, onlyMeta, lowMeta, paths,
		)
		fmt.Println()

	} else if options.InferIdPath != nil {
		result, err := fsext.InferIdsFromWorkPath(*options.InferIdPath)
		logext.MaybeFatal(err, "cannot infer work id from pattern %v", *options.InferIdPath)
		if len(*result) == 0 {
			logext.Warning("no ids could be inferred from pattern %v", *options.InferIdPath)
			return
		}

		if options.Path == nil {
			q := queue.FromMap(result, kind, size, onlyMeta)
			d.ScheduleQueue(q)
		} else {
			paths := []string{path}
			q := queue.FromMapWithPaths(result, kind, size, onlyMeta, paths)
			d.ScheduleQueue(q)
		}

	} else if options.QueuePath != nil {
		paths := []string{path}
		q, warnings, err := fsext.ReadQueue(*options.QueuePath, kind, size, onlyMeta, paths)
		logext.MaybeWarnings(warnings, "while reading the list from %v", *options.QueuePath)
		logext.MaybeFatal(err, "cannot read the list from %v", *options.QueuePath)
		if len(*q) == 0 {
			logext.Warning("no works found in the list %v", *options.QueuePath)
			return
		}

		d.ScheduleQueue(q)
	}

	fmt.Println(d)

	logext.EnableProgress()
	defer logext.DisableProgress()

	logext.Info("download started")
	d.Run() // TODO: confirmation by user (or -y flag)
	d.WaitDone()
	logext.Info("download finished")
}

func configSession(password *string) (*config.Storage, *string) {
	storage, err := config.Open(password)
	if err != nil && password != nil {
		logext.Fatal("cannot open config storage: %v", err)
		panic("unreachable")
	} else if err != nil {
		logext.Warning("cannot open config storage: %v", err)
		promptDefaultConfig()
		return storage, nil
	}

	if sessionId, err := storage.SessionId(); err != nil && password != nil {
		logext.Fatal("cannot read session id: %v", err)
		panic("unreachable")
	} else if err != nil {
		if newStorage, sessionId := promptPassword(); newStorage != nil {
			return newStorage, sessionId
		} else {
			return storage, sessionId
		}
	} else if sessionId == nil && password != nil {
		logext.Fatal("no session id were configured, but password was provided")
		panic("unreachable")
	} else if sessionId == nil {
		logext.Info("no session id were configured, using only anonymous requests")
		return storage, nil
	} else {
		return storage, sessionId
	}
}

func promptPassword() (*config.Storage, *string) {
	for tries := 0; ; tries++ {
		password, err := passwordPrompt.Run()
		if err != nil {
			logext.Warning("failed to read password: %v", err)
			promptNoAuthorization()
			return nil, nil
		}

		storage, err := config.Open(&password)
		if err != nil {
			logext.Warning("cannot open config storage: %v", err)
			promptNoAuthorization()
			return nil, nil
		}

		if sessionId, err := storage.SessionId(); err == nil && sessionId != nil {
			return storage, sessionId
		} else if err == nil {
			logext.Info("no session id were configured, using only anonymous requests")
			return storage, nil
		} else if tries == 2 {
			logext.Warning("cannot read session id: %v", err)
			promptNoAuthorization()
			return storage, nil
		}
	}
}

func promptDefaultConfig() {
	_, option, err := deafultConfigPrompt.Run()
	logext.MaybeFatal(err, "failed to read the choice")
	if err != nil || option != YesOption {
		os.Exit(1)
	}
}

func promptNoAuthorization() {
	_, option, err := noAuthorizationPrompt.Run()
	logext.MaybeFatal(err, "failed to read the choice")
	if err != nil || option != YesOption {
		os.Exit(1)
	}
}
