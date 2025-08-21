package download

import (
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

// TODO: look through interractive mode, add descriptions and explanations

var idOption = "Download by ID"
var myPublicBookmarksOption = "Download my public bookmarks"
var myPrivateBookmarksOption = "Download my private bookmarks"
var userBookmarksOption = "Download bookmarks of other user"
var inferIdOption = "Infer IDs from path"
var listOption = "Download from list"

var sourceSelect = promptui.Select{
	Label: "Download mode",
	Items: []string{
		idOption, myPublicBookmarksOption, myPrivateBookmarksOption,
		userBookmarksOption, inferIdOption, listOption,
	},
}

var idsPrompt = promptui.Prompt{
	Label: "Work IDs",
	Validate: func(input string) error {
		_, err := parseIds(input)
		return err
	},
}

var userIdPrompt = promptui.Prompt{
	Label:    "User ID",
	Validate: utils.ValidateNumber("user ID must be a number"),
}

var inferIdsPrompt = promptui.Prompt{
	Label: "Paths to directories or patterns to infer IDs from (comma-separated)",
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

var listsPrompt = promptui.Prompt{
	Label: "Paths to download list files (comma-separated)",
}

var artworkOption = "Artwork"
var novelOption = "Novel"

func kindSelect(withLists bool) *promptui.Select {
	const label = "Type of work to download"
	const withListsLabel = "Default type of work to download"

	return &promptui.Select{
		Label: utils.If(withLists, withListsLabel, label),
		Items: []string{artworkOption, novelOption},
	}
}

var tagsPrompt = promptui.Prompt{
	Label: "User-assigned tags (comma-separated or empty)",
}

var rangePrompt = promptui.Prompt{
	Label: "Download range (from:to, from:, :to or empty)",
	Validate: func(input string) error {
		_, _, err := parseRange(input)
		return err
	},
}

var downloadAllOption = "Download metadata and images"
var downloadMetaOption = "Only download metadata"

func onlyMetaSelect(withLists bool) *promptui.Select {
	const label = "Only download metadata"
	const withListsLabel = "Only download metadata by default"

	return &promptui.Select{
		Label: utils.If(withLists, withListsLabel, label),
		Items: []string{downloadAllOption, downloadMetaOption},
	}
}

var lowMetaOption = "Save partial metadata"
var fullMetaOption = "Get full metadata"

var lowMetaSelect = promptui.Select{
	Label: "Don't get full metadata (less requests)",
	Items: []string{lowMetaOption, fullMetaOption},
}

var skipsPrompt = promptui.Prompt{
	Label: "Skip works present in the directory (path, infer id pattern or nothing)",
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

var untilSkipOption = "Fetch new bookmarks until fully downloaded bookmark page is reached"
var allPagesOption = "Fetch and check all bookmarks"

var untilSkipSelect = promptui.Select{
	Label: "Pick bookmark pages fetching strategy",
	Items: []string{untilSkipOption, allPagesOption},
}

var thumbnailSizeOption = "Thumbnail"
var smallSizeOption = "Small"
var mediumSizeOption = "Medium"
var originalSizeOption = "Original"

func sizeSelect(withLists bool) *promptui.Select {
	const label = "Size of downloaded images"
	const withListsLabel = "Default size of downloaded images"

	return &promptui.Select{
		Label:     utils.If(withLists, withListsLabel, label),
		Items:     []string{thumbnailSizeOption, smallSizeOption, mediumSizeOption, originalSizeOption},
		CursorPos: 3,
	}
}

var inferredPathOption = "Save to inferred path"
var customPathOption = "Specify different path"

var pathSelect = promptui.Select{
	Label: "Where to save downloaded works?",
	Items: []string{inferredPathOption, customPathOption},
}

func pathsPrompt(withLists bool) *promptui.Prompt {
	const label = "Save to directory (one or multiple comma-separated)"
	const withListsLabel = "Default saving path (one or multiple comma-separated)"

	return &promptui.Prompt{
		Label:    utils.If(withLists, withListsLabel, label),
		Validate: fsext.WorkPathPatternValid,
	}
}

var rulesPrompt = promptui.Prompt{
	Label: "Path to download rules file (leave empty to download all)",
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

var ignoreRulesPrompt = promptui.Select{
	Label: "Ignore global download rules?",
	Items: []string{YesOption, NoOption},
}

var ignoreLimitsPrompt = promptui.Select{
	Label: "Use default limits?",
	Items: []string{YesOption, NoOption},
}
