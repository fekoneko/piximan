package dto

import (
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

type Work struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	UserId        string `json:"userId"`
	UserName      string `json:"userName"`
	XRestrict     uint8  `json:"xRestrict"`
	AiType        uint8  `json:"aiType"`
	IsOriginal    bool   `json:"isOriginal"`
	PageCount     uint64 `json:"pageCount"`
	ViewCount     uint64 `json:"viewCount"`
	BookmarkCount uint64 `json:"bookmarkCount"`
	LikeCount     uint64 `json:"likeCount"`
	CommentCount  uint64 `json:"commentCount"`
	UploadDate    string `json:"uploadDate"`
	SeriesNavData struct {
		SeriesId *uint64 `json:"seriesId"`
		Order    *uint64 `json:"order"`
		Title    *string `json:"title"`
	} `json:"seriesNavData"`
	Tags struct {
		Tags [](struct {
			Tag string `json:"tag"`
		}) `json:"tags"`
	} `json:"tags"`
}

func (dto *Work) FromDto(kind work.Kind, downloadTime time.Time) *work.Work {
	id, _ := strconv.ParseUint(dto.Id, 10, 64)
	userId, _ := strconv.ParseUint(dto.UserId, 10, 64)

	uploadTime, err := time.Parse(time.RFC3339, dto.UploadDate)
	localUploadTime := uploadTime.Local()
	localUploadTimePtr := &localUploadTime
	if err != nil {
		localUploadTimePtr = nil
	}

	tags := make([]string, len(dto.Tags.Tags))
	for i, tag := range dto.Tags.Tags {
		tags[i] = tag.Tag
	}

	return &work.Work{
		Id:           id,
		Title:        dto.Title,
		Kind:         kind,
		Description:  formatDescription(dto.Description),
		UserId:       userId,
		UserName:     dto.UserName,
		Restriction:  work.RestrictionFromUint(dto.XRestrict),
		AiKind:       work.AiKindFromUint(dto.AiType),
		IsOriginal:   dto.IsOriginal,
		NumPages:     dto.PageCount,
		NumViews:     dto.ViewCount,
		NumBookmarks: dto.BookmarkCount,
		NumLikes:     dto.LikeCount,
		NumComments:  dto.CommentCount,
		UploadTime:   localUploadTimePtr,
		DownloadTime: downloadTime.Local(),
		SeriesId:     dto.SeriesNavData.SeriesId,
		SeriesTitle:  dto.SeriesNavData.Title,
		SeriesOrder:  dto.SeriesNavData.Order,
		Tags:         tags,
	}
}

// This function converts the work description with HTML to plain text. It does the following:
// - replaces <br> tags with \n
// - replaces <a> tags with their href attribute values
// - removes all other HTML tags
// - unescapes HTML entities
func formatDescription(description string) string {
	tagStart := 0
	textStart := 0
	ignoreText := false
	builder := strings.Builder{}

	for i, b := range description {
		switch b {
		case '<':
			tagStart = i + 1
			if !ignoreText {
				builder.WriteString(description[textStart:i])
			}
		case '>':
			textStart = i + 1
			tag := strings.ToLower(description[tagStart:i])
			if strings.HasPrefix(tag, "br ") || tag == "br/" || tag == "br" {
				builder.WriteString("\n")
			} else if strings.HasPrefix(tag, "a ") {
				ignoreText = true
				hrefStart := strings.Index(tag, "href=") + len("href=") + 1
				if hrefStart != len("href=") && hrefStart < len(tag) {
					hrefLength := strings.IndexAny(tag[hrefStart:], "\"'")
					if hrefLength != -1 {
						href := tag[hrefStart : hrefStart+hrefLength]
						builder.WriteString(href)
					}
				}
			} else if tag == "/a" || strings.HasPrefix(tag, "/a ") {
				ignoreText = false
			}
		}
	}

	if textStart < len(description) && !ignoreText {
		builder.WriteString(description[textStart:])
	}

	return html.UnescapeString(builder.String())
}
