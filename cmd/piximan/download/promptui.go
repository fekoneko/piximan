package download

import (
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

var sourceSelectLabel = "Download mode"
var idOption = "Download by ID"
var myPublicBookmarksOption = "Download my public bookmarks"
var myPrivateBookmarksOption = "Download my private bookmarks"
var userBookmarksOption = "Download bookmarks of other user"
var inferIdOption = "Infer IDs from path"
var queueOption = "Download from list"

var sourceSelect = promptui.Select{
	Label: sourceSelectLabel,
	Items: []string{
		idOption, myPublicBookmarksOption, myPrivateBookmarksOption,
		userBookmarksOption, inferIdOption, queueOption,
	},
}

var idPromptLabel = "Work IDs"

var idPrompt = promptui.Prompt{
	Label: idPromptLabel,
	Validate: func(input string) error {
		_, err := parseIds(input)
		return err
	},
}

var userIdPromptLabel = "User ID"

var userIdPrompt = promptui.Prompt{
	Label:    userIdPromptLabel,
	Validate: utils.ValidateNumber("user ID must be a number"),
}

var inferIdPromptLabel = "Path pattern"

var inferIdPrompt = promptui.Prompt{
	Label:    inferIdPromptLabel,
	Validate: fsext.InferIdPatternValid,
}

var listPromptLabel = "Path to download list file"

var listPrompt = promptui.Prompt{
	Label: listPromptLabel,
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

var tagsPromptLabel = "User-assigned tags (comma-separated or empty)"

var tagsPrompt = promptui.Prompt{
	Label: tagsPromptLabel,
}

var rangePromptLabel = "Download range (from:to, from:, :to or empty)"

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

var skipPromptLabel = "Skip works present in the directory (path, infer id pattern or nothing)"

var skipPrompt = promptui.Prompt{
	Label: skipPromptLabel,
	Validate: func(input string) error {
		if fsext.IsInferIdPattern(input) {
			return fsext.InferIdPatternValid(input)
		}
		return nil
	},
}

var untilSkipSelectLabel = "Pick bookmark pages fetching strategy"
var untilSkipOption = "Fetch new bookmarks until fully downloaded bookmark page is reached"
var allPagesOption = "Fetch and check all bookmarks"

var untilSkipSelect = promptui.Select{
	Label: untilSkipSelectLabel,
	Items: []string{untilSkipOption, allPagesOption},
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
		Label:    utils.If(withQueue, pathPromptWithQueueLabel, pathPromptLabel),
		Validate: fsext.WorkPathPatternValid,
	}
}

var rulesPromptLabel = "Path to download rules file (leave empty to download all)"

var rulesPrompt = promptui.Prompt{
	Label: rulesPromptLabel,
}

var passwordPrompt = promptui.Prompt{
	Label:       "Password",
	Mask:        '*',
	HideEntered: true,
}

var YesOption = "Yes"
var NoOption = "No"

var deafultConfigPrompt = promptui.Select{
	Label: "Use default config and anonymous requests?",
	Items: []string{"Yes", "No"},
}

var noAuthorizationPrompt = promptui.Select{
	Label: "Use only anonymous requests?",
	Items: []string{"Yes", "No"},
}
