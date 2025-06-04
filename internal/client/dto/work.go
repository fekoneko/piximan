package dto

import (
	"html"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

type Work struct {
	Id            *string `json:"id"`
	Title         *string `json:"title"`
	Description   *string `json:"description"`
	UserId        *string `json:"userId"`
	UserName      *string `json:"userName"`
	XRestrict     *uint8  `json:"xRestrict"`
	AiType        *uint8  `json:"aiType"`
	IsOriginal    *bool   `json:"isOriginal"`
	PageCount     *uint64 `json:"pageCount"`
	ViewCount     *uint64 `json:"viewCount"`
	BookmarkCount *uint64 `json:"bookmarkCount"`
	LikeCount     *uint64 `json:"likeCount"`
	CommentCount  *uint64 `json:"commentCount"`
	UploadDate    *string `json:"uploadDate"`
	SeriesNavData struct {
		SeriesId *uint64 `json:"seriesId"`
		Order    *uint64 `json:"order"`
		Title    *string `json:"title"`
	} `json:"seriesNavData"`
	Tags struct {
		Tags *[](struct {
			Tag *string `json:"tag"`
		}) `json:"tags"`
	} `json:"tags"`
}

func (dto *Work) FromDto(kind *work.Kind, downloadTime time.Time) *work.Work {
	var tags *[]string
	if dto.Tags.Tags != nil {
		tags = utils.ToPtr(make([]string, len(*dto.Tags.Tags)))
		for i, tag := range *dto.Tags.Tags {
			if tag.Tag == nil {
				tags = nil
				break
			}
			(*tags)[i] = *tag.Tag
		}
	}

	return &work.Work{
		Id:           utils.ParseUint64Ptr(dto.Id),
		Title:        dto.Title,
		Kind:         kind,
		Description:  formatDescription(dto.Description),
		UserId:       utils.ParseUint64Ptr(dto.UserId),
		UserName:     dto.UserName,
		Restriction:  utils.MapPtr(dto.XRestrict, work.RestrictionFromUint),
		AiKind:       utils.MapPtr(dto.XRestrict, work.AiKindFromUint),
		Original:     dto.IsOriginal,
		NumPages:     dto.PageCount,
		NumViews:     dto.ViewCount,
		NumBookmarks: dto.BookmarkCount,
		NumLikes:     dto.LikeCount,
		NumComments:  dto.CommentCount,
		UploadTime:   utils.ParseLocalTimePtr(dto.UploadDate),
		DownloadTime: utils.ToPtr(downloadTime.Local()),
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
func formatDescription(description *string) *string {
	if description == nil {
		return nil
	}

	tagStart := 0
	textStart := 0
	builder := strings.Builder{}

	for i, b := range *description {
		switch b {
		case '<':
			tagStart = i + 1
			builder.WriteString((*description)[textStart:i])
		case '>':
			textStart = i + 1
			tag := strings.ToLower((*description)[tagStart:i])
			if strings.HasPrefix(tag, "br ") || tag == "br/" || tag == "br" {
				builder.WriteString("\n")
			}
		}
	}

	if textStart < len(*description) {
		builder.WriteString((*description)[textStart:])
	}

	result := html.UnescapeString(builder.String())
	return &result
}
