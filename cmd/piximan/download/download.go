package download

import (
	"fmt"
	"strconv"

	"github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/downloader"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/termext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

func download(options *options) {
	size := utils.FromPtrTransform(options.Size, image.SizeFromUint, image.SizeDefault)
	kind := utils.FromPtrTransform(options.Kind, queue.ItemKindFromString, queue.ItemKindDefault)
	private := utils.FromPtr(options.Private, false)
	onlyMeta := utils.FromPtr(options.OnlyMeta, false)
	lowMeta := utils.FromPtr(options.LowMeta, false)
	path := utils.FromPtr(options.Path, "")

	d := chooseDownloader(options.Password)

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

	logext.EnableRequestSlots()
	defer logext.DisableRequestSlots()

	logext.Info("download started")
	d.Run() // TODO: confirmation by user (or -y flag)
	d.WaitDone()
	logext.Info("download finished")
}

func chooseDownloader(passwordPtr *string) *downloader.Downloader {
	password := utils.FromPtr(passwordPtr, "")

	storage, err := config.Open(password)
	if err != nil && passwordPtr != nil {
		logext.Fatal("cannot open session id storage: %v", err)
		panic("unreachable")
	} else if err != nil {
		logext.Warning("cannot open session id storage, using only anonymous requests: %v\n", err)
		return downloader.New(nil)
	}

	if err := storage.Read(); err != nil && passwordPtr != nil {
		logext.Fatal("cannot read session id: %v", err)
		panic("unreachable")
	} else if err != nil {
		return promptPassword()
	} else if storage.SessionId == nil && passwordPtr != nil {
		logext.Fatal("no session id were configured, but password was provided")
		panic("unreachable")
	} else if storage.SessionId == nil {
		logext.Info("no session id were configured, using only anonymous requests")
		return downloader.New(nil)
	} else {
		return downloader.New(storage.SessionId)
	}
}

var passwordPrompt = promptui.Prompt{
	Label:       "Password",
	Mask:        '*',
	HideEntered: true,
}

func promptPassword() *downloader.Downloader {
	for tries := 0; ; tries++ {
		password, err := passwordPrompt.Run()
		if err != nil {
			logext.Warning("failed to read password, using only anonymous requests: %v\n", err)
			return downloader.New(nil)
		}

		storage, err := config.Open(password)
		if err != nil {
			logext.Warning("cannot open session id storage, using only anonymous requests: %v\n", err)
			return downloader.New(nil)
		}

		if err = storage.Read(); err == nil && storage.SessionId != nil {
			return downloader.New(storage.SessionId)
		} else if err == nil {
			logext.Info("no session id were configured, using only anonymous requests\n")
			return downloader.New(nil)
		} else if tries == 2 {
			logext.Warning("cannot read session id, using only anonymous requests: %v\n", err)
			return downloader.New(nil)
		}
	}
}
