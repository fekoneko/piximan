package download

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/utils"
)

func interactive() {
	ids, bookmarks, inferIdPath, queuePath := selectSource()
	withQueue := queuePath != nil
	withInferId := inferIdPath != nil
	withBookmarks := bookmarks != nil

	kind := selectKind(withQueue)
	fromOffset, toOffset := promptRange(withBookmarks)
	onlyMeta := selectOnlyMeta(withQueue)
	lowMeta := selectLowMeta(withBookmarks)
	size := selectSize(withQueue, onlyMeta)
	path := promptPath(withInferId, withQueue)

	fmt.Println()
	download(&options{
		Ids:         ids,
		Bookmarks:   bookmarks,
		QueuePath:   queuePath,
		InferIdPath: inferIdPath,
		Kind:        &kind,
		Size:        size,
		OnlyMeta:    &onlyMeta,
		LowMeta:     lowMeta,
		FromOffset:  fromOffset,
		ToOffset:    toOffset,
		Path:        path,
	})
}

func selectSource() (*[]uint64, *string, *string, *string) {
	_, mode, err := sourceSelect.Run()
	logext.MaybeFatal(err, "failed to read mode")

	switch mode {
	case idOption:
		idsString, err := idPrompt.Run()
		logext.MaybeFatal(err, "failed to read IDs")
		ids, err := parseIds(idsString)
		logext.MaybeFatal(err, "failed to parse IDs")
		return &ids, nil, nil, nil

	case myBookmarksOption:
		return nil, utils.ToPtr("my"), nil, nil

	case userBookmarksOption:
		userId, err := userIdPrompt.Run()
		logext.MaybeFatal(err, "failed to read user ID")
		return nil, utils.ToPtr(userId), nil, nil

	case inferIdOption:
		inferIdPath, err := inferIdPathPrompt.Run()
		logext.MaybeFatal(err, "failed to read pattern")
		return nil, nil, &inferIdPath, nil

	case queueOption:
		queuePath, err := queuePathPrompt.Run()
		logext.MaybeFatal(err, "failed to read list path")
		return nil, nil, nil, &queuePath

	default:
		logext.Fatal("incorrect download mode: %v", mode)
		panic("unreachable")
	}
}

func selectKind(withQueue bool) string {
	_, kind, err := kindSelect(withQueue).Run()
	logext.MaybeFatal(err, "failed to read work type")

	switch kind {
	case artworkOption:
		return queue.ItemKindArtworkString
	case novelOption:
		return queue.ItemKindNovelString
	default:
		logext.Fatal("invalid worktype: %s", kind)
		panic("unreachable")
	}
}

func promptRange(withBookmarks bool) (*uint64, *uint64) {
	if !withBookmarks {
		return nil, nil
	}

	rangeString, err := rangePrompt.Run()
	logext.MaybeFatal(err, "failed to read range")
	fromOffset, toOffset, err := parseRange(rangeString)
	logext.MaybeFatal(err, "failed to parse range")

	return fromOffset, toOffset
}

func selectOnlyMeta(withQueue bool) bool {
	_, option, err := onlyMetaSelect(withQueue).Run()
	logext.MaybeFatal(err, "failed to read downloaded files choice")

	switch option {
	case downloadAllOption:
		return false
	case downloadMetaOption:
		return true
	default:
		logext.Fatal("incorrect downloaded files choice: %v", option)
		panic("unreachable")
	}
}

func selectLowMeta(withBookmarks bool) *bool {
	if !withBookmarks {
		return nil
	}

	_, option, err := lowMetaSelect.Run()
	logext.MaybeFatal(err, "failed to read low metadata choice")

	switch option {
	case lowMetaOption:
		return utils.ToPtr(true)
	case fullMetaOption:
		return utils.ToPtr(false)
	default:
		logext.Fatal("incorrect low metadata choice: %v", option)
		panic("unreachable")
	}
}

func selectSize(withQueue bool, onlyMeta bool) *uint {
	if onlyMeta {
		return nil
	}

	_, size, err := sizeSelect(withQueue).Run()
	logext.MaybeFatal(err, "failed to read size")

	switch size {
	case thumbnailSizeOption:
		result := uint(image.SizeThumbnail)
		return &result
	case smallSizeOption:
		result := uint(image.SizeSmall)
		return &result
	case mediumSizeOption:
		result := uint(image.SizeMedium)
		return &result
	case originalSizeOption:
		result := uint(image.SizeOriginal)
		return &result
	default:
		logext.Fatal("incorrect size: %v", size)
		panic("unreachable")
	}
}

func promptPath(withInferId bool, withQueue bool) *string {
	if withInferId {
		_, pathChoice, err := pathSelect.Run()
		logext.MaybeFatal(err, "failed to read path choice")
		if pathChoice == inferredPathOption {
			return nil
		}
	}

	path, err := pathPrompt(withQueue).Run()
	logext.MaybeFatal(err, "failed to read path")
	return &path
}
