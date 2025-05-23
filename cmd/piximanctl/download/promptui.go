package download

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

var sourceSelectLabel = "Download mode"
var idOption = "Download by ID"
var myBookmarksOption = "Download my bookmarks"
var userBookmarksOption = "Download bookmarks of other user"
var inferIdOption = "Infer IDs from path"
var queueOption = "Download from list"

var sourceSelect = promptui.Select{
	Label: sourceSelectLabel,
	Items: []string{idOption, myBookmarksOption, userBookmarksOption, inferIdOption, queueOption},
}

var idPromptLabel = "Work IDs"

var idPrompt = promptui.Prompt{
	Label: idPromptLabel,
	Validate: func(input string) error {
		if _, err := parseIds(input); err != nil {
			return fmt.Errorf("IDs must be a comma-separated list of numbers")
		}
		return nil
	},
}

var userIdPromptLabel = "User ID"

var userIdPrompt = promptui.Prompt{
	Label: userIdPromptLabel,
	Validate: func(input string) error {
		if _, err := strconv.ParseUint(input, 10, 64); err != nil {
			return fmt.Errorf("user ID must be a number")
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

var rangePromptLabel = "Download range (from:to, from:, :to)"

var rangePrompt = promptui.Prompt{
	Label: rangePromptLabel,
	Validate: func(input string) error {
		_, _, err := parseRange(input)
		return err
	},
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

var lowMetaSelectLabel = "Don't get full metadata (less requests)"
var lowMetaOption = "Save partial metadata"
var fullMetaOption = "Get full metadata"

var lowMetaSelect = promptui.Select{
	Label: lowMetaSelectLabel,
	Items: []string{lowMetaOption, fullMetaOption},
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
