package download

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/utils"
)

func interactive() {
	ids, bookmarks, private, inferId, list := selectSource()
	withQueue := list != nil
	withInferId := inferId != nil
	withBookmarks := bookmarks != nil

	kind := selectKind(withQueue)
	tags := promptTags(withBookmarks)
	fromOffset, toOffset := promptRange(withBookmarks)
	onlyMeta := selectOnlyMeta(withQueue)
	lowMeta := selectLowMeta(withBookmarks, kind, onlyMeta)
	skip := promptSkip(withBookmarks)
	withSkip := skip != nil
	untilSkip := selectUntilSkip(withSkip)
	size := selectSize(withQueue, onlyMeta)
	path := promptPath(withInferId, withQueue)
	rules := promptRules()

	fmt.Println()
	download(&options{
		Ids:        ids,
		Bookmarks:  bookmarks,
		List:       list,
		InferId:    inferId,
		Kind:       &kind,
		Size:       size,
		OnlyMeta:   &onlyMeta,
		Rules:      rules,
		Skip:       skip,
		Tags:       tags,
		FromOffset: fromOffset,
		ToOffset:   toOffset,
		Private:    private,
		LowMeta:    lowMeta,
		UntilSkip:  untilSkip,
		Path:       path,
	})
}

func selectSource() (ids *[]uint64, bookmarks *string, private *bool, inferId *[]string, list *string) {
	_, mode, err := sourceSelect.Run()
	logger.MaybeFatal(err, "failed to read mode")

	switch mode {
	case idOption:
		idsString, err := idPrompt.Run()
		logger.MaybeFatal(err, "failed to read IDs")
		parsed, err := parseIds(idsString)
		logger.MaybeFatal(err, "failed to parse IDs")
		ids = &parsed

	case myPublicBookmarksOption:
		bookmarks = utils.ToPtr("my")
		private = utils.ToPtr(false)

	case myPrivateBookmarksOption:
		bookmarks = utils.ToPtr("my")
		private = utils.ToPtr(true)

	case userBookmarksOption:
		userId, err := userIdPrompt.Run()
		logger.MaybeFatal(err, "failed to read user ID")
		bookmarks = utils.ToPtr(userId)

	case inferIdOption:
		result, err := inferIdPrompt.Run()
		logger.MaybeFatal(err, "failed to read infer id option")
		parsed := parseStrings(result)
		inferId = &parsed

	case queueOption:
		result, err := listPrompt.Run()
		logger.MaybeFatal(err, "failed to read list path")
		list = &result

	default:
		logger.Fatal("incorrect download mode: %v", mode)
	}
	return
}

func selectKind(withQueue bool) string {
	_, kind, err := kindSelect(withQueue).Run()
	logger.MaybeFatal(err, "failed to read work type")

	switch kind {
	case artworkOption:
		return queue.ItemKindArtworkString
	case novelOption:
		return queue.ItemKindNovelString
	default:
		logger.Fatal("invalid worktype: %s", kind)
		panic("unreachable")
	}
}

func promptTags(withBookmarks bool) *[]string {
	if !withBookmarks {
		return nil
	}

	tagsString, err := tagsPrompt.Run()
	logger.MaybeFatal(err, "failed to read tags")
	tags := parseStrings(tagsString)
	if len(tags) == 0 {
		return nil
	}
	return &tags
}

func promptRange(withBookmarks bool) (fromOffset *uint64, toOffset *uint64) {
	if !withBookmarks {
		return
	}

	rangeString, err := rangePrompt.Run()
	logger.MaybeFatal(err, "failed to read range")
	fromOffset, toOffset, err = parseRange(rangeString)
	logger.MaybeFatal(err, "failed to parse range")
	return
}

func selectOnlyMeta(withQueue bool) bool {
	_, option, err := onlyMetaSelect(withQueue).Run()
	logger.MaybeFatal(err, "failed to read downloaded files choice")

	switch option {
	case downloadAllOption:
		return false
	case downloadMetaOption:
		return true
	default:
		logger.Fatal("incorrect downloaded files choice: %v", option)
		panic("unreachable")
	}
}

func selectLowMeta(withBookmarks bool, kind string, onlyMeta bool) *bool {
	if !withBookmarks || (kind == queue.ItemKindNovelString && !onlyMeta) {
		return nil
	}

	_, option, err := lowMetaSelect.Run()
	logger.MaybeFatal(err, "failed to read low metadata choice")

	switch option {
	case lowMetaOption:
		return utils.ToPtr(true)
	case fullMetaOption:
		return utils.ToPtr(false)
	default:
		logger.Fatal("incorrect low metadata choice: %v", option)
		panic("unreachable")
	}
}

func promptSkip(withBookmarks bool) *[]string {
	if !withBookmarks {
		return nil
	}
	skipString, err := skipPrompt.Run()
	logger.MaybeFatal(err, "failed to read skip option")
	skip := parseStrings(skipString)
	if len(skip) == 0 {
		return nil
	}
	return &skip
}

func selectUntilSkip(withSkip bool) *bool {
	if !withSkip {
		return nil
	}
	_, option, err := untilSkipSelect.Run()
	logger.MaybeFatal(err, "failed to read until skip flag")

	switch option {
	case untilSkipOption:
		return utils.ToPtr(true)
	case allPagesOption:
		return utils.ToPtr(false)
	default:
		logger.Fatal("incorrect until skip flag choice: %v", option)
		panic("unreachable")
	}
}

func selectSize(withQueue bool, onlyMeta bool) *uint {
	if onlyMeta {
		return nil
	}

	_, size, err := sizeSelect(withQueue).Run()
	logger.MaybeFatal(err, "failed to read size")

	switch size {
	case thumbnailSizeOption:
		result := uint(imageext.SizeThumbnail)
		return &result
	case smallSizeOption:
		result := uint(imageext.SizeSmall)
		return &result
	case mediumSizeOption:
		result := uint(imageext.SizeMedium)
		return &result
	case originalSizeOption:
		result := uint(imageext.SizeOriginal)
		return &result
	default:
		logger.Fatal("incorrect size: %v", size)
		panic("unreachable")
	}
}

func promptPath(withInferId bool, withQueue bool) *string {
	if withInferId {
		_, pathChoice, err := pathSelect.Run()
		logger.MaybeFatal(err, "failed to read path choice")
		if pathChoice == inferredPathOption {
			return nil
		}
	}

	path, err := pathPrompt(withQueue).Run()
	logger.MaybeFatal(err, "failed to read path")
	return &path
}

func promptRules() *string {
	rules, err := rulesPrompt.Run()
	logger.MaybeFatal(err, "failed to read rules path")
	if rules == "" {
		return nil
	}
	return &rules
}
