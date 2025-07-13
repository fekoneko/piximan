package dto

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

type Novel struct {
	Work
	Content            *string `json:"content"`
	CoverUrl           *string `json:"coverUrl"`
	TextEmbeddedImages map[string]struct {
		Urls struct {
			Thumb    *string `json:"128x128"`
			Small    *string `json:"480mw"`
			Regular  *string `json:"1200x1200"`
			Original *string `json:"original"`
		} `json:"urls"`
	} `json:"textEmbeddedImages"`
}

// TODO: wrap lines at 60 - 80 characters on word boundaries

// Provided size is only used to determine embedded image urls.
// If you don't need novel content and images, pass nil instead.
func (dto *Novel) FromDto(downloadTime time.Time, size *imageext.Size) (
	w *work.Work, coverUrl *string, uploadedImages NovelUpladedImages,
	pixivImages NovelPixivImages, pages NovelPages, withPages bool,
) {
	w = dto.Work.FromDto(utils.ToPtr(work.KindNovel), downloadTime)
	coverUrl = dto.CoverUrl

	if size == nil || dto.Content == nil {
		return w, coverUrl, nil, nil, nil, false
	}

	matches := contentRegexp.FindAllStringSubmatchIndex(*dto.Content, -1)
	pixivImages = make(NovelPixivImages)
	uploadedImages = make(NovelUpladedImages)
	imageIndexes := make([]int, 0)

	for _, match := range matches {
		if match[uploadedImage] >= 0 {
			if dto.TextEmbeddedImages == nil {
				return w, coverUrl, nil, nil, nil, false
			}
			idString := (*dto.Content)[match[uploadedImageId]:match[uploadedImageId+1]]
			url, ok := dto.uploadedImageUrl(idString, *size)
			if !ok {
				return w, coverUrl, nil, nil, nil, false
			}
			index, ok := utils.MapFindValue(uploadedImages, url)
			if !ok {
				index = len(uploadedImages) + 1
				uploadedImages[index] = url
			}
			imageIndexes = append(imageIndexes, index)

		} else if match[pixivImage] >= 0 {
			idString := (*dto.Content)[match[pixivImageId]:match[pixivImageId+1]]
			id, _ := strconv.ParseUint(idString, 10, 64)
			index, ok := utils.MapFindValue(pixivImages, id)
			if !ok {
				index = len(pixivImages) + 1
				pixivImages[index] = id
			}
			imageIndexes = append(imageIndexes, index)
		}
	}

	pages = func(
		imageName func(index int) string,
		pageName func(index int) string,
	) []fsext.Asset {
		return finishParsingContent(dto.Content, matches, imageIndexes, imageName, pageName)
	}

	return w, coverUrl, uploadedImages, pixivImages, pages, true
}

func (dto *Novel) uploadedImageUrl(idString string, size imageext.Size) (url string, ok bool) {
	if dto.TextEmbeddedImages == nil {
		return "", false
	} else if urls, ok := dto.TextEmbeddedImages[idString]; !ok {
		return "", false
	} else if size == imageext.SizeThumbnail && urls.Urls.Thumb != nil {
		return *urls.Urls.Thumb, true
	} else if size == imageext.SizeSmall && urls.Urls.Small != nil {
		return *urls.Urls.Small, true
	} else if size == imageext.SizeMedium && urls.Urls.Regular != nil {
		return *urls.Urls.Regular, true
	} else if size == imageext.SizeOriginal && urls.Urls.Original != nil {
		return *urls.Urls.Original, true
	} else {
		return "", false
	}
}

// Convert novel content from pixiv format to markdown. This does the following:
// - 2 or more \n -> \n\n
// - \n -> <br>
// - [newpage] -> write to next page and trim empty lines at the beginning and at the end
// - [uploadedimage:{id}] -> ![Illustration]({name})
// - [pixivimage:{id}] -> ![Illustration]({name})
// - [[rb:{word} > {ruby}]] -> <ruby>{word}<rt>{ruby}</rt></ruby>
// - [chapter:{title}] -> # {title}
// - [jump:{page}] -> [{page}]({name})
// - [[jumpuri:{title} > {url}]] -> [{title}]({url})
func finishParsingContent(
	content *string, matches [][]int, imageIndexes []int,
	imageName func(index int) string, pageName func(index int) string,
) []fsext.Asset {
	assets := make([]fsext.Asset, 0, 1)
	builder := strings.Builder{}
	prevEnd := 0
	pageNumber := 1
	numImages := 0

	for _, match := range matches {
		if prevEnd < match[0] {
			builder.WriteString((*content)[prevEnd:match[0]])
		}
		prevEnd = match[1]

		if match[newPage] >= 0 {
			builder.WriteString("\n")
			asset := fsext.Asset{Bytes: []byte(builder.String()), Name: pageName(pageNumber)}
			assets = append(assets, asset)
			builder.Reset()
			pageNumber++

		} else if match[startNewLines] >= 0 {
		} else if match[endNewLines] >= 0 {

		} else if match[newParagraph] >= 0 {
			builder.WriteString("\n\n")

		} else if match[newLine] >= 0 {
			builder.WriteString("<br>")

		} else if match[uploadedImage] >= 0 {
			builder.WriteString("![Illustration](<")
			builder.WriteString(imageName(imageIndexes[numImages]))
			builder.WriteString(">)")
			numImages++

		} else if match[pixivImage] >= 0 {
			builder.WriteString("![Illustration](<")
			builder.WriteString(imageName(imageIndexes[numImages]))
			builder.WriteString(">)")
			numImages++

		} else if match[ruby] >= 0 {
			word := (*content)[match[rubyWord]:match[rubyWord+1]]
			ruby := (*content)[match[rubyRuby]:match[rubyRuby+1]]
			builder.WriteString("<ruby>")
			builder.WriteString(word)
			builder.WriteString("<rt>")
			builder.WriteString(ruby)
			builder.WriteString("</rt></ruby>")

		} else if match[title] >= 0 {
			text := (*content)[match[titleText]:match[titleText+1]]
			builder.WriteString("# ")
			builder.WriteString(text)

		} else if match[pageLink] >= 0 {
			pageString := (*content)[match[pageLinkPage]:match[pageLinkPage+1]]
			page, _ := strconv.ParseInt(pageString, 10, 64)
			builder.WriteByte('[')
			builder.WriteString(pageString)
			builder.WriteString("](<")
			builder.WriteString(pageName(int(page)))
			builder.WriteString(">)")

		} else if match[urlLink] >= 0 {
			text := (*content)[match[urlLinkText]:match[urlLinkText+1]]
			url := (*content)[match[urlLinkUrl]:match[urlLinkUrl+1]]
			builder.WriteByte('[')
			builder.WriteString(text)
			builder.WriteString("](")
			builder.WriteString(url)
			builder.WriteByte(')')

		} else if match[specialChar] >= 0 {
			char := (*content)[match[specialChar]:match[specialChar+1]]
			builder.WriteByte('\\')
			builder.WriteString(char)
		}
	}

	if prevEnd < len(*content) {
		builder.WriteString((*content)[prevEnd:len(*content)])
	}
	builder.WriteString("\n")
	asset := fsext.Asset{Bytes: []byte(builder.String()), Name: pageName(pageNumber)}
	assets = append(assets, asset)

	return assets
}

var contentRegexp = regexp.MustCompile(
	`([ 　\n\t]*\[newpage\][ 　\n\t]*)|` +
		`(^[ 　\n\t]*)|` +
		`([ 　\n\t]*$)|` +
		`([ 　\t]*\n[ 　\t]*\n[ 　\n\t]*)|` +
		`([ 　\t]*\n[ 　\t]*)|` +
		`(\[uploadedimage:([0-9]+)\])|` +
		`(\[pixivimage:([0-9]+)\])|` +
		`(\[\[rb:(.+) *> *(.+)\]\])|` +
		`(\[chapter:(.+)\])|` +
		`(\[jump:([0-9]+)\])|` +
		`(\[\[jumpuri:(.+) *> *(.+)\]\])|` +
		`([#&*+\-<[\\_|~` + "`" + `])`,
)

const (
	newPage         = 2 * 1
	startNewLines   = 2 * 2
	endNewLines     = 2 * 3
	newParagraph    = 2 * 4
	newLine         = 2 * 5
	uploadedImage   = 2 * 6
	uploadedImageId = 2 * 7
	pixivImage      = 2 * 8
	pixivImageId    = 2 * 9
	ruby            = 2 * 10
	rubyWord        = 2 * 11
	rubyRuby        = 2 * 12
	title           = 2 * 13
	titleText       = 2 * 14
	pageLink        = 2 * 15
	pageLinkPage    = 2 * 16
	urlLink         = 2 * 17
	urlLinkText     = 2 * 18
	urlLinkUrl      = 2 * 19
	specialChar     = 2 * 20
)

// Map: index -> URL, used for downloading novel illustrations.
type NovelUpladedImages map[int]string

// Map: index -> Pixiv work ID, used for downloading novel illustrations.
type NovelPixivImages map[int]uint64

// Call to finish parsing novel contents when all asset names are known.
type NovelPages = func(
	imageName func(index int) string,
	pageName func(page int) string,
) []fsext.Asset
