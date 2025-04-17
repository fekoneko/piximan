package download

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
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
	if err != nil {
		fmt.Printf("failed to read work type: %v\n", err)
		os.Exit(1)
	}
	switch kind {
	case ArtworkOption:
		return queue.ItemKindArtworkString
	case NovelOption:
		return queue.ItemKindNovelString
	default:
		fmt.Printf("invalid worktype: %s\n", kind)
		os.Exit(1)
	}
	return ""
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
	if err != nil {
		fmt.Printf("failed to read mode: %v\n", err)
		os.Exit(1)
	}
	switch mode {
	case idModeOption:
		idString, err := idPrompt.Run()
		if err != nil {
			fmt.Printf("failed to read ID: %v\n", err)
			os.Exit(1)
		}
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Printf("failed to parse ID: %v\n", err)
			os.Exit(1)
		}
		return &id, nil
	case inferIdModeOption:
		inferId, err := inferIdPrompt.Run()
		if err != nil {
			fmt.Printf("failed to read pattern: %v\n", err)
			os.Exit(1)
		}
		return nil, &inferId
	default:
		fmt.Printf("incorrect download mode: %v\n", mode)
		os.Exit(1)
	}
	return nil, nil
}

var downloadAllOption = "Download metadata and images"
var downloadMetaOption = "Only download metadata"
var downloadFilesSelect = promptui.Select{
	Label: "Downloaded files",
	Items: []string{downloadAllOption, downloadMetaOption},
}

func selectOnlyMeta() bool {
	_, downloadFiles, err := downloadFilesSelect.Run()
	if err != nil {
		fmt.Printf("failed to read downloaded files choice: %v\n", err)
		os.Exit(1)
	}
	switch downloadFiles {
	case downloadAllOption:
		return false
	case downloadMetaOption:
		return true
	default:
		fmt.Printf("incorrect downloaded files choice: %v\n", downloadFiles)
		os.Exit(1)
	}
	return false
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
	if err != nil {
		fmt.Printf("failed to read size: %v\n", err)
		os.Exit(1)
	}
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
		fmt.Printf("incorrect size: %v\n", size)
		os.Exit(1)
	}
	return uint(image.SizeDefault)
}

var inferredPathOption = "Save to inferred path"
var customPathOption = "Specify different path"
var pathSelect = promptui.Select{
	Label: "Where to save downloaded works?",
	Items: []string{inferredPathOption, customPathOption},
}

func selectAskPath() bool {
	_, pathChoice, err := pathSelect.Run()
	if err != nil {
		fmt.Printf("failed to read path choice: %v\n", err)
		os.Exit(1)
	}
	return pathChoice == customPathOption
}

var pathPrompt = promptui.Prompt{
	Label: "Save to directory",
}

func promptPath() string {
	path, err := pathPrompt.Run()
	if err != nil {
		fmt.Printf("failed to read path: %v\n", err)
		os.Exit(1)
	}
	return path
}
