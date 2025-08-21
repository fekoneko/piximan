package download

import (
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

// TODO: look through interractive mode, add descriptions and explanations

var sourceSelectLabel = "Download mode"
var idOption = "Download by ID"
var myPublicBookmarksOption = "Download my public bookmarks"
var myPrivateBookmarksOption = "Download my private bookmarks"
var userBookmarksOption = "Download bookmarks of other user"
var inferIdOption = "Infer IDs from path"
var listOption = "Download from list"

var sourceSelect = promptui.Select{
	Label: sourceSelectLabel,
	Items: []string{
		idOption, myPublicBookmarksOption, myPrivateBookmarksOption,
		userBookmarksOption, inferIdOption, listOption,
	},
}

var idsPromptLabel = "Work IDs"

var idsPrompt = promptui.Prompt{
	Label: idsPromptLabel,
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

var inferIdsPromptLabel = "Paths to directories or patterns to infer IDs from (comma-separated)"

var inferIdsPrompt = promptui.Prompt{
	Label: inferIdsPromptLabel,
	Validate: func(input string) error {
		for _, s := range parseStrings(input) {
			if !fsext.IsInferIdPattern(input) {
				continue
			}
			if err := fsext.InferIdPatternValid(s); err != nil {
				return err
			}
		}
		return nil
	},
}

var listsPromptLabel = "Paths to download list files (comma-separated)"

var listsPrompt = promptui.Prompt{
	Label: listsPromptLabel,
}

var kindSelectLabel = "Type of work to download"
var kindSelectwithListsLabel = "Default type of work to download"
var artworkOption = "Artwork"
var novelOption = "Novel"

func kindSelect(withLists bool) *promptui.Select {
	return &promptui.Select{
		Label: utils.If(withLists, kindSelectwithListsLabel, kindSelectLabel),
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
var onlyMetaSelectwithListsLabel = "Only download metadata by default"
var downloadAllOption = "Download metadata and images"
var downloadMetaOption = "Only download metadata"

func onlyMetaSelect(withLists bool) *promptui.Select {
	return &promptui.Select{
		Label: utils.If(withLists, onlyMetaSelectwithListsLabel, onlyMetaSelectLabel),
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

var skipsPromptLabel = "Skip works present in the directory (path, infer id pattern or nothing)"

var skipsPrompt = promptui.Prompt{
	Label: skipsPromptLabel,
	Validate: func(input string) error {
		for _, s := range parseStrings(input) {
			if !fsext.IsInferIdPattern(input) {
				continue
			}
			if err := fsext.InferIdPatternValid(s); err != nil {
				return err
			}
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
var sizeSelectwithListsLabel = "Default size of downloaded images"
var thumbnailSizeOption = "Thumbnail"
var smallSizeOption = "Small"
var mediumSizeOption = "Medium"
var originalSizeOption = "Original"

func sizeSelect(withLists bool) *promptui.Select {
	return &promptui.Select{
		Label:     utils.If(withLists, sizeSelectwithListsLabel, sizeSelectLabel),
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

var pathsPromptLabel = "Save to directory (one or multiple comma-separated)"
var pathsPromptwithListsLabel = "Default saving path (one or multiple comma-separated)"

func pathsPrompt(withLists bool) *promptui.Prompt {
	return &promptui.Prompt{
		Label:    utils.If(withLists, pathsPromptwithListsLabel, pathsPromptLabel),
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

var ignoreAuthorizationPrompt = promptui.Select{
	Label: "Use only anonymous requests?",
	Items: []string{YesOption, NoOption},
}

var ignoreLimitsPrompt = promptui.Select{
	Label: "Use default limits?",
	Items: []string{YesOption, NoOption},
}
