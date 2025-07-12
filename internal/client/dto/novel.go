package dto

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
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

func (dto *Novel) FromDto(downloadTime time.Time) (w *work.Work, pages *[]string, coverUrl *string) {
	w = dto.Work.FromDto(utils.ToPtr(work.KindNovel), downloadTime)
	pages = parseContent(dto.Content)

	return w, pages, dto.CoverUrl
}

// TODO: download embedded images
// TODO: wrap lines at 60 - 80 characters on word boundaries

// Convert novel content from pixiv format to markdown. This does the following:
// - 2 or more \n -> \n\n
// - \n -> <br>
// - [newpage] -> write to next page and trim empty lines at the beginning and at the end
// - [uploadedimage:{id}] -> ![Illustration](./{id}. {title}.{ext})
// - [pixivimage:{id}] -> ![Illustration](./{id}. {title}.{ext})
// - [[rb:{word} > {ruby}]] -> <ruby>{word}<rt>{ruby}</rt></ruby>
// - [chapter:{title}] -> # {title}
// - [jump:{page}] -> [{page}](./{page}. {title}.md)
// - [[jumpuri:{title} > {url}]] -> [{title}]({url})
func parseContent(content *string) (pages *[]string) {
	if content == nil {
		return nil
	}

	pages = utils.ToPtr(make([]string, 0, 1))
	matches := contentRegexp.FindAllStringSubmatchIndex(*content, -1)
	builder := strings.Builder{}
	prevEnd := 0

	for _, match := range matches {
		if prevEnd < match[0] {
			builder.WriteString((*content)[prevEnd:match[0]])
		}
		prevEnd = match[1]

		if match[newPage] >= 0 {
			builder.WriteString("\n")
			*pages = append(*pages, builder.String())
			builder.Reset()
		} else if match[startNewLines] >= 0 {
		} else if match[endNewLines] >= 0 {
		} else if match[newParagraph] >= 0 {
			builder.WriteString("\n\n")
		} else if match[newLine] >= 0 {
			builder.WriteString("<br>")
		} else if match[uploadedImage] >= 0 {
			builder.WriteString("![Illustration]()") // TODO: filename
		} else if match[pixivImage] >= 0 {
			builder.WriteString("![Illustration]()") // TODO: filename
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
			page, _ := strconv.ParseUint(pageString, 10, 64)
			builder.WriteByte('[')
			builder.WriteString(pageString)
			builder.WriteString("](")
			builder.WriteString(fsext.NovelPageAssetName(page))
			builder.WriteByte(')')
		} else if match[urlLink] >= 0 {
			text := (*content)[match[urlLinkText]:match[urlLinkText+1]]
			url := (*content)[match[urlLinkUrl]:match[urlLinkUrl+1]]
			builder.WriteByte('[')
			builder.WriteString(text)
			builder.WriteString("](")
			builder.WriteString(url)
			builder.WriteByte(')')
		}
	}

	if prevEnd < len(*content) {
		builder.WriteString((*content)[prevEnd:len(*content)])
	}
	builder.WriteString("\n")
	*pages = append(*pages, builder.String())

	return pages
}

var contentRegexp = regexp.MustCompile(
	`(\n*\[newpage\]\n*)|` +
		`(^\n*)|` +
		`(\n*$)|` +
		`(\n{2,})|` +
		`(\n)|` +
		`(\[uploadedimage:([0-9]+)\])|` +
		`(\[pixivimage:([0-9]+)\])|` +
		`(\[\[rb:(.+) *> *(.+)\]\])|` +
		`(\[chapter:(.+)\])|` +
		`(\[jump:([0-9]+)\])|` +
		`(\[\[jumpuri:(.+) *> *(.+)\]\])`,
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
)
