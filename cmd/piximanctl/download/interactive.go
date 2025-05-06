package download

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/utils"
	"github.com/manifoldco/promptui"
)

func interactive() {
	ids, inferIdPath, queuePath := selectMode()
	withQueue := queuePath != nil
	withInferId := inferIdPath != nil

	kind := selectKind(withQueue)
	onlyMeta := selectOnlyMeta(withQueue)
	size := selectSize(withQueue, onlyMeta)
	path := promptPath(withInferId, withQueue)

	fmt.Println()
	download(&options{
		Ids:         ids,
		Kind:        &kind,
		Size:        size,
		Path:        path,
		InferIdPath: inferIdPath,
		OnlyMeta:    &onlyMeta,
	})
}

func selectMode() (*[]uint64, *string, *string) {
	_, mode, err := modeSelect.Run()
	logext.MaybeFatal(err, "failed to read mode")

	switch mode {
	case idModeOption:
		idsString, err := idPrompt.Run()
		logext.MaybeFatal(err, "failed to read IDs")
		ids, err := parseIdsString(idsString)
		logext.MaybeFatal(err, "failed to parse IDs")
		return &ids, nil, nil

	case inferIdModeOption:
		inferIdPath, err := inferIdPathPrompt.Run()
		logext.MaybeFatal(err, "failed to read pattern")
		return nil, &inferIdPath, nil

	case queueModeOption:
		queuePath, err := queuePathPrompt.Run()
		logext.MaybeFatal(err, "failed to read list path")
		return nil, nil, &queuePath

	default:
		logext.Fatal("incorrect download mode: %v", mode)
	}
	panic("unreachable")
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
	}
	panic("unreachable")
}

func selectOnlyMeta(withQueue bool) bool {
	_, downloadFiles, err := onlyMetaSelect(withQueue).Run()
	logext.MaybeFatal(err, "failed to read downloaded files choice")

	switch downloadFiles {
	case downloadAllOption:
		return false
	case downloadMetaOption:
		return true
	default:
		logext.Fatal("incorrect downloaded files choice: %v", downloadFiles)
	}
	panic("unreachable")
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
	}
	panic("unreachable")
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

func parseIdsString(idsString string) ([]uint64, error) {
	idSubstrs := strings.Split(idsString, ",")
	ids := []uint64{}

	for _, idSubstr := range idSubstrs {
		trimmed := strings.TrimSpace(idSubstr)
		if trimmed == "" {
			continue
		}
		id, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no IDs provided")
	}

	return ids, nil
}

var modeSelectLabel = "Download mode"
var idModeOption = "Download by ID"
var inferIdModeOption = "Infer IDs from path"
var queueModeOption = "Download from list"

var modeSelect = promptui.Select{
	Label: modeSelectLabel,
	Items: []string{idModeOption, inferIdModeOption, queueModeOption},
}

var idPromptLabel = "Work IDs"

var idPrompt = promptui.Prompt{
	Label: idPromptLabel,
	Validate: func(input string) error {
		if _, err := parseIdsString(input); err != nil {
			return fmt.Errorf("IDs must be a comma-separated list of numbers")
		}
		return nil
	},
}

var inferIdPathPromptLabel = "Path pattern"

var inferIdPathPrompt = promptui.Prompt{
	Label: inferIdPathPromptLabel,
	Validate: func(input string) error {
		if !strings.Contains(input, "{id}") {
			return fmt.Errorf("pattern must contain {id}")
		}
		return nil
	},
}

var queuePathPromptLabel = "Path to YAML list"

var queuePathPrompt = promptui.Prompt{
	Label: queuePathPromptLabel,
}

var kindSelectLabel = "Type of work to download"
var kindSelectWithQueueLabel = "Default type of work to download"
var artworkOption = "Artwork"
var novelOption = "Novel"

func kindSelect(withQueue bool) *promptui.Select {
	return &promptui.Select{
		Label: utils.If(withQueue, kindSelectWithQueueLabel, kindSelectLabel),
		Items: []string{artworkOption, novelOption},
	}
}

var onlyMetaSelectLabel = "Only download metadata"
var onlyMetaSelectWithQueueLabel = "Only download metadata by default"
var downloadAllOption = "Download metadata and images"
var downloadMetaOption = "Only download metadata"

func onlyMetaSelect(withQueue bool) *promptui.Select {
	return &promptui.Select{
		Label: utils.If(withQueue, onlyMetaSelectWithQueueLabel, onlyMetaSelectLabel),
		Items: []string{downloadAllOption, downloadMetaOption},
	}
}

var sizeSelectLabel = "Size of downloaded images"
var sizeSelectWithQueueLabel = "Default size of downloaded images"
var thumbnailSizeOption = "Thumbnail"
var smallSizeOption = "Small"
var mediumSizeOption = "Medium"
var originalSizeOption = "Original"

func sizeSelect(withQueue bool) *promptui.Select {
	return &promptui.Select{
		Label:     utils.If(withQueue, sizeSelectWithQueueLabel, sizeSelectLabel),
		Items:     []string{thumbnailSizeOption, smallSizeOption, mediumSizeOption, originalSizeOption},
		CursorPos: 3,
	}
}

var pathSelectLabel = "Where to save downloaded works?"
var inferredPathOption = "Save to inferred path"
var customPathOption = "Specify different path"

var pathSelect = promptui.Select{
	Label: pathSelectLabel,
	Items: []string{inferredPathOption, customPathOption},
}

var pathPromptLabel = "Save to directory"
var pathPromptWithQueueLabel = "Default saving path"

func pathPrompt(withQueue bool) *promptui.Prompt {
	return &promptui.Prompt{
		Label: utils.If(withQueue, pathPromptWithQueueLabel, pathPromptLabel),
	}
}
