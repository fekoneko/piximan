package download

import (
	"fmt"
	"strings"

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

var tagPromptLabel = "User-assigned tag (leave empty for any)"

var tagPrompt = promptui.Prompt{
	Label: tagPromptLabel,
}

var bookmarksConstraintSelectLabel = "Limit downloaded bookmarks"
var withRangeOption = "By offset range"
var withTimeOption = "By time the work was bookmarked"
var noBookmarksConstraintOption = "Download all bookmarks"

var bookmarksConstraintSelect = promptui.Select{
	Label: bookmarksConstraintSelectLabel,
	Items: []string{withRangeOption, withTimeOption},
}

var rangePromptLabel = "Download range (from:to, from:, :to)"

var rangePrompt = promptui.Prompt{
	Label: rangePromptLabel,
	Validate: func(input string) error {
		_, _, err := parseRange(input)
		return err
	},
}

var olderThanPromptLabel = "Bookmarked before (YYYY-MM-DD or empty)"

var olderThanPrompt = promptui.Prompt{
	Label: olderThanPromptLabel,
	Validate: func(input string) error {
		_, err := parseTime(input)
		return err
	},
}

var newerThanPromptLabel = "Bookmarked after (YYYY-MM-DD or empty)"

var newerThanPrompt = promptui.Prompt{
	Label: newerThanPromptLabel,
	Validate: func(input string) error {
		_, err := parseTime(input)
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
