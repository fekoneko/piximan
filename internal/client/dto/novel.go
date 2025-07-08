package dto

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
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

func (dto *Novel) FromDto(downloadTime time.Time) (w *work.Work, content *string, coverUrl *string) {
	w = dto.Work.FromDto(utils.ToPtr(work.KindNovel), downloadTime)

	content, _ = parseContent(w, dto.Content)

	return w, content, dto.CoverUrl
}

// TODO: parse page breaks
// TODO: parse page titles
// TODO: parse page links
// TODO: parse url links
// TODO: parse ruby (furigana)
// TODO: ensure image extensions are correct when quolity is 'original'

// Convert novel content from pixiv format to markdown. This does the following:
// - doubles the line breaks
// - replaces embedded image annotations with the appropriate markdown links
// Also returns a map to later translate Pixiv image ids to local IDs.
func parseContent(w *work.Work, content *string) (_ *string, imageMap map[uint64]uint64) {
	imageMap = make(map[uint64]uint64)
	charsOnLine := 0
	latestImageId := 1
	builder := strings.Builder{}

	if content == nil {
		return nil, imageMap
	}

	for i := 0; i < len(*content); i++ {
		if (*content)[i] == '\n' {
			builder.WriteString("\n\n")
			charsOnLine = 0
		} else if id := embeddedImage(content, &i); id != 0 {
			latestImageId++
			title := utils.FromPtr(w.Title, "unknown")
			s := fmt.Sprintf("![Embedded Image](./%03d. %v.jpg)", latestImageId, title)
			builder.WriteString(s)
		} else {
			builder.WriteByte((*content)[i])
			charsOnLine++
		}
	}

	parsed := builder.String()
	return &parsed, imageMap
}

// Determine if embedded image annotation starts at index i. If it is, return the image ID,
// otherwise return 0. Content is assumed to be not nil and have a length of at least i+1.
// The func tion advances i to the last character of found annotation.
func embeddedImage(content *string, i *int) uint64 {
	const annotation = "uploadedimage:"

	if (*content)[*i] != '[' ||
		len(*content) <= *i+len(annotation)+3 ||
		(*content)[*i+1:*i+len(annotation)+1] != annotation {
		return 0
	}

	idStart := *i + len(annotation) + 1
	idLength := strings.IndexByte((*content)[idStart:], ']')
	if idLength == -1 {
		return 0
	}

	idEnd := idStart + idLength
	id := (*content)[idStart:idEnd]
	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0
	}

	*i = idEnd
	return parsed
}
