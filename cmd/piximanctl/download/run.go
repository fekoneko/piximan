package download

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/manifoldco/promptui"
)

type flags struct {
	id       *uint64
	kind     *string
	size     *uint
	path     *string
	inferId  *string
	onlyMeta *bool
}

func Run() {
	flags := flags{
		id:       flag.Uint64("id", 0, ""),
		kind:     flag.String("type", "artwork", ""),
		size:     flag.Uint("size", uint(image.SizeDefault), ""),
		path:     flag.String("path", "", ""),
		inferId:  flag.String("inferid", "", ""),
		onlyMeta: flag.Bool("onlymeta", false, ""),
	}
	flag.Usage = help.RunDownload
	flag.Parse()

	if flag.NArg() != 0 {
		flagext.BadUsage("too many arguments")
	}

	if flag.NFlag() == 0 {
		interactive(flags)
	} else {
		nonInteractive(flags)
	}
}

// TODO: write descriptions
var ArtworkOption = "Artwork"
var NovelOption = "Novel"
var kindSelect = promptui.Select{
	Label: "Type of work to download",
	Items: []string{ArtworkOption, NovelOption},
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
var downloadAllOption = "Download metadata and images"
var downloadMetaOption = "Only download metadata"
var downloadFilesSelect = promptui.Select{
	Label: "Downloaded files",
	Items: []string{downloadAllOption, downloadMetaOption},
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
var inferredPathOption = "Save to inferred path"
var customPathOption = "Specify different path"
var pathSelect = promptui.Select{
	Label: "Where to save downloaded works?",
	Items: []string{inferredPathOption, customPathOption},
}
var pathPrompt = promptui.Prompt{
	Label: "Save to directory",
}

func interactive(flags flags) {
	_, kind, err := kindSelect.Run()
	if err != nil {
		fmt.Printf("failed to read work type: %v\n", err)
		os.Exit(1)
	}
	switch kind {
	case ArtworkOption:
		*flags.kind = queue.ItemKindArtworkString
	case NovelOption:
		*flags.kind = queue.ItemKindNovelString
	default:
		fmt.Printf("invalid worktype: %s\n", kind)
		os.Exit(1)
	}

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
		*flags.id = id
	case inferIdModeOption:
		inferId, err := inferIdPrompt.Run()
		if err != nil {
			fmt.Printf("failed to read pattern: %v\n", err)
			os.Exit(1)
		}
		*flags.inferId = inferId
	default:
		fmt.Printf("incorrect download mode: %v\n", mode)
		os.Exit(1)
	}

	_, downloadFiles, err := downloadFilesSelect.Run()
	if err != nil {
		fmt.Printf("failed to read downloaded files choice: %v\n", err)
		os.Exit(1)
	}
	switch downloadFiles {
	case downloadAllOption:
		*flags.onlyMeta = false
	case downloadMetaOption:
		*flags.onlyMeta = true
	default:
		fmt.Printf("incorrect downloaded files choice: %v\n", downloadFiles)
		os.Exit(1)
	}

	if downloadFiles == downloadAllOption {
		_, size, err := sizeSelect.Run()
		if err != nil {
			fmt.Printf("failed to read size: %v\n", err)
			os.Exit(1)
		}
		switch size {
		case thumbnailSizeOption:
			*flags.size = uint(image.SizeThumbnail)
		case smallSizeOption:
			*flags.size = uint(image.SizeSmall)
		case mediumSizeOption:
			*flags.size = uint(image.SizeMedium)
		case originalSizeOption:
			*flags.size = uint(image.SizeOriginal)
		default:
			fmt.Printf("incorrect size: %v\n", size)
			os.Exit(1)
		}
	}

	askPath := true
	if mode == inferIdModeOption {
		_, pathChoice, err := pathSelect.Run()
		if err != nil {
			fmt.Printf("failed to read path choice: %v\n", err)
			os.Exit(1)
		}
		if pathChoice == inferredPathOption {
			askPath = false
		}
	}

	if askPath {
		path, err := pathPrompt.Run()
		if err != nil {
			fmt.Printf("failed to read path: %v\n", err)
			os.Exit(1)
		}
		*flags.path = path
	}

	fmt.Println()
	continueDownload(flags, mode == inferIdModeOption, askPath)
}

func nonInteractive(flags flags) {
	if !flagext.Provided("id") && !flagext.Provided("inferid") {
		flagext.BadUsage("one of these arguments is not provided: -id, -inferid")
	}
	if flagext.Provided("id") && flagext.Provided("inferid") {
		flagext.BadUsage("providing these arguments together is not supporded: -id, -inferid")
	}
	if flagext.Provided("type") && flagext.Provided("size") && *flags.kind == queue.ItemKindNovelString {
		flagext.BadUsage("cannot use -size argument with -type novel")
	}
	if flagext.Provided("type") && queue.ValidItemKindString(*flags.kind) {
		flagext.BadUsage("invalid argument value: -type")
	}
	if flagext.Provided("size") && *flags.size > 3 {
		flagext.BadUsage("invalid argument value: -size")
	}
	continueDownload(flags, flagext.Provided("inferid"), flagext.Provided("path"))
}

// TODO: pass not flags here
func continueDownload(flags flags, isInferIdProvided bool, isPathProvided bool) {
	d := downloader.New()
	size := image.SizeFromUint(*flags.size)
	kind := queue.ItemKindFromString(*flags.kind)

	if isInferIdProvided {
		result, err := pathext.InferIdsFromWorkPath(*flags.inferId)
		if err != nil {
			fmt.Printf("cannot infer work id from pattern %v: %v\n", *flags.inferId, err)
			os.Exit(1)
		}
		q := queue.FromMap(result, kind, size, *flags.onlyMeta)
		if isPathProvided {
			for i := range *q {
				(*q)[i].Paths = []string{*flags.path}
			}
		}
		fmt.Print(q, "\n\n")
		d.ScheduleQueue(q)
	} else {
		paths := []string{*flags.path}
		d.Schedule(*flags.id, kind, size, *flags.onlyMeta, paths)
	}

	for d.Listen() != nil {
	}
}
