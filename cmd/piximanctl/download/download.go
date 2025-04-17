package download

import (
	"fmt"

	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
	"github.com/fekoneko/piximan/pkg/util"
	"github.com/manifoldco/promptui"
)

func download(options *options) {
	size := util.FromPtrTransform(options.Size, image.SizeFromUint, image.SizeDefault)
	kind := util.FromPtrTransform(options.King, queue.ItemKindFromString, queue.ItemKindDefault)
	onlyMeta := util.FromPtr(options.OnlyMeta, false)
	path := util.FromPtr(options.Path, "")

	d := chooseDownloader(options.Password)

	if options.InferId != nil {
		result, err := pathext.InferIdsFromWorkPath(*options.InferId)
		logext.MaybeFatal(err, "cannot infer work id from pattern %v", *options.InferId)

		q := queue.FromMap(result, kind, size, onlyMeta)
		if options.Path != nil {
			for i := range *q {
				(*q)[i].Paths = []string{path}
			}
		}
		fmt.Print(q, "\n\n")
		d.ScheduleQueue(q)
	} else {
		d.Schedule(*options.Id, kind, size, onlyMeta, []string{path})
	}

	for d.Listen() != nil {
	}
}

func chooseDownloader(passwordPtr *string) *downloader.Downloader {
	password := util.FromPtr(passwordPtr, "")

	storage, err := secretstorage.Open(password)
	if err != nil && passwordPtr != nil {
		logext.Fatal("cannot open session id storage: %v", err)
	} else if err != nil {
		logext.Warning("cannot open session id storage, using only anonymous requests: %v\n", err)
		return downloader.New()
	}

	if err := storage.Read(); err != nil && passwordPtr != nil {
		logext.Fatal("cannot read session id: %v", err)
	} else if err != nil {
		return promptPassword()
	} else if storage.SessionId == nil && passwordPtr != nil {
		logext.Fatal("no session id were configured, but password was provided")
	} else if storage.SessionId == nil {
		logext.Info("no session id were configured, using only anonymous requests")
		return downloader.New()
	} else {
		return downloader.NewAuthorized(*storage.SessionId)
	}
	panic("unreachable")
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
			return downloader.New()
		}

		storage, err := secretstorage.Open(password)
		if err != nil {
			logext.Warning("cannot open session id storage, using only anonymous requests: %v\n", err)
			return downloader.New()
		}

		if err = storage.Read(); err == nil && storage.SessionId != nil {
			return downloader.NewAuthorized(*storage.SessionId)
		} else if err == nil {
			logext.Info("no session id were configured, using only anonymous requests\n")
			return downloader.New()
		} else if tries == 2 {
			logext.Warning("cannot read session id, using only anonymous requests: %v\n", err)
			return downloader.New()
		}
	}
}
