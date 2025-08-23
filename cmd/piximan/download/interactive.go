package download

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/utils"
)

func interactive() {
	ids, bookmarks, private, inferIds, lists := selectSource()
	withLists := lists != nil
	withInferIds := inferIds != nil
	withBookmarks := bookmarks != nil
	kind := selectKind(withLists)
	tags := promptTags(withBookmarks)
	fromOffset, toOffset := promptRange(withBookmarks)
	onlyMeta := selectOnlyMeta(withLists)
	lowMeta := selectLowMeta(withBookmarks, kind, onlyMeta)
	skips := promptSkips(withBookmarks)
	withSkips := skips != nil
	untilSkip := selectUntilSkip(withSkips)
	size := selectSize(withLists, onlyMeta)
	language := selectLanguage()
	paths := promptPaths(withInferIds, withLists)
	rules := promptRules()

	fmt.Println()
	download(&options{
		Ids:        ids,
		Bookmarks:  bookmarks,
		Lists:      lists,
		InferIds:   inferIds,
		Kind:       &kind,
		Size:       size,
		Language:   language,
		OnlyMeta:   &onlyMeta,
		Rules:      rules,
		Skips:      skips,
		Tags:       tags,
		FromOffset: fromOffset,
		ToOffset:   toOffset,
		Private:    private,
		LowMeta:    lowMeta,
		UntilSkip:  untilSkip,
		Paths:      paths,
	})
}

func selectSource() (ids *[]uint64, bookmarks *string, private *bool, inferIds, lists *[]string) {
	_, mode, err := sourceSelect.Run()
	logger.MaybeFatal(err, "failed to read mode")

	switch mode {
	case idOption:
		idsString, err := idsPrompt.Run()
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
		result, err := inferIdsPrompt.Run()
		logger.MaybeFatal(err, "failed to read infer id option")
		parsed := parseStrings(result)
		inferIds = &parsed

	case listOption:
		result, err := listsPrompt.Run()
		logger.MaybeFatal(err, "failed to read list paths")
		parsed := parseStrings(result)
		lists = &parsed

	default:
		logger.Fatal("incorrect download mode: %v", mode)
	}
	return
}

func selectKind(withLists bool) string {
	_, kind, err := kindSelect(withLists).Run()
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

func selectOnlyMeta(withLists bool) bool {
	_, option, err := onlyMetaSelect(withLists).Run()
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

func promptSkips(withBookmarks bool) *[]string {
	if !withBookmarks {
		return nil
	}
	skipsString, err := skipsPrompt.Run()
	logger.MaybeFatal(err, "failed to read skip option")
	skips := parseStrings(skipsString)
	if len(skips) == 0 {
		return nil
	}
	return &skips
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

func selectSize(withLists, onlyMeta bool) *uint {
	if onlyMeta {
		return nil
	}

	_, size, err := sizeSelect(withLists).Run()
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

func selectLanguage() *string {
	_, language, err := languageSelect(work.LanguageDefault).Run() // TODO: default language from config
	logger.MaybeFatal(err, "failed to read language choice")

	switch language {
	case japaneseOption:
		return utils.ToPtr(work.LanguageJapaneseString)
	case englishOption:
		return utils.ToPtr(work.LanguageEnglishString)
	default:
		logger.Fatal("incorrect language: %v", language)
		panic("unreachable")
	}
}

func promptPaths(withInferIds, withLists bool) *[]string {
	if withInferIds {
		_, pathChoice, err := pathSelect.Run()
		logger.MaybeFatal(err, "failed to read path choice")
		if pathChoice == inferredPathOption {
			return nil
		}
	}

	pathsString, err := pathsPrompt(withLists).Run()
	logger.MaybeFatal(err, "failed to read paths")
	paths := parseStrings(pathsString)
	return &paths
}

func promptRules() *[]string {
	rules, err := rulesPrompt.Run()
	logger.MaybeFatal(err, "failed to read rules paths")
	parsed := parseStrings(rules)
	if len(parsed) == 0 {
		return nil
	}
	return &parsed
}
