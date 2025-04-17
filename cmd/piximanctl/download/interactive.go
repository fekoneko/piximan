package download

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/manifoldco/promptui"
)

func interactive() {
	kind := selectKind()
	id, inferId := selectMode()
	onlyMeta := selectOnlyMeta()

	var size *uint
	if !onlyMeta {
		s := selectSize()
		size = &s
	}

	var path *string
	if inferId == nil || selectAskPath() {
		p := promptPath()
		path = &p
	}

	fmt.Println()
	download(&options{
		Id:       id,
		King:     &kind,
		Size:     size,
		Path:     path,
		InferId:  inferId,
		OnlyMeta: &onlyMeta,
	})
}

var ArtworkOption = "Artwork"
var NovelOption = "Novel"
var kindSelect = promptui.Select{
	Label: "Type of work to download",
	Items: []string{ArtworkOption, NovelOption},
}

func selectKind() string {
	_, kind, err := kindSelect.Run()
	logext.MaybeFatal(err, "failed to read work type")

	switch kind {
	case ArtworkOption:
		return queue.ItemKindArtworkString
	case NovelOption:
		return queue.ItemKindNovelString
	default:
		logext.Fatal("invalid worktype: %s", kind)
	}
	panic("unreachable")
}

var idModeOption = "Download by ID"
var inferIdModeOption = "Infer IDs from path"
var modeSelect = promptui.Select{
	Label: "Download mode",
	Items: []string{idModeOption, inferIdModeOption},
}
var idPrompt = promptui.Prompt{
	Label: "Work ID",
	Validate: func(input string) error {
		if _, err := strconv.ParseUint(input, 10, 64); err != nil {
			return fmt.Errorf("ID must be a number")
		}
		return nil
	},
}
var inferIdPrompt = promptui.Prompt{
	Label: "Path pattern",
	Validate: func(input string) error {
		if !strings.Contains(input, "{id}") {
			return fmt.Errorf("pattern must contain {id}")
		}
		return nil
	},
}

func selectMode() (*uint64, *string) {
	_, mode, err := modeSelect.Run()
	logext.MaybeFatal(err, "failed to read mode")

	switch mode {
	case idModeOption:
		idString, err := idPrompt.Run()
		logext.MaybeFatal(err, "failed to read ID")

		id, err := strconv.ParseUint(idString, 10, 64)
		logext.MaybeFatal(err, "failed to parse ID")

		return &id, nil
	case inferIdModeOption:
		inferId, err := inferIdPrompt.Run()
		logext.MaybeFatal(err, "failed to read pattern")

		return nil, &inferId
	default:
		logext.Fatal("incorrect download mode: %v", mode)
	}
	panic("unreachable")
}

var downloadAllOption = "Download metadata and images"
var downloadMetaOption = "Only download metadata"
var downloadFilesSelect = promptui.Select{
	Label: "Downloaded files",
	Items: []string{downloadAllOption, downloadMetaOption},
}

func selectOnlyMeta() bool {
	_, downloadFiles, err := downloadFilesSelect.Run()
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

var thumbnailSizeOption = "Thumbnail"
var smallSizeOption = "Small"
var mediumSizeOption = "Medium"
var originalSizeOption = "Original"
var sizeSelect = promptui.Select{
	Label:     "Size of downloaded images",
	Items:     []string{thumbnailSizeOption, smallSizeOption, mediumSizeOption, originalSizeOption},
	CursorPos: 3,
}

func selectSize() uint {
	_, size, err := sizeSelect.Run()
	logext.MaybeFatal(err, "failed to read size")

	switch size {
	case thumbnailSizeOption:
		return uint(image.SizeThumbnail)
	case smallSizeOption:
		return uint(image.SizeSmall)
	case mediumSizeOption:
		return uint(image.SizeMedium)
	case originalSizeOption:
		return uint(image.SizeOriginal)
	default:
		logext.Fatal("incorrect size: %v", size)
	}
	panic("unreachable")
}

var inferredPathOption = "Save to inferred path"
var customPathOption = "Specify different path"
var pathSelect = promptui.Select{
	Label: "Where to save downloaded works?",
	Items: []string{inferredPathOption, customPathOption},
}

func selectAskPath() bool {
	_, pathChoice, err := pathSelect.Run()
	logext.MaybeFatal(err, "failed to read path choice")
	return pathChoice == customPathOption
}

var pathPrompt = promptui.Prompt{
	Label: "Save to directory",
}

func promptPath() string {
	path, err := pathPrompt.Run()
	logext.MaybeFatal(err, "failed to read path")
	return path
}
