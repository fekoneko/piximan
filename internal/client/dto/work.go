package dto

import (
	"html"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

type Work struct {
	Id            *string `json:"id"`
	Title         *string `json:"title"`       // TODO: translation
	Description   *string `json:"description"` // TODO: translation
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
	CreateDate    *string `json:"createDate"`
	SeriesNavData struct {
		SeriesId any     `json:"seriesId"`
		Order    *uint64 `json:"order"`
		Title    *string `json:"title"`
	} `json:"seriesNavData"`
	Tags struct {
		Tags *[](struct {
			Tag *string `json:"tag"` // TODO: translation
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

	var seriesId *uint64
	seriesIdType := reflect.TypeOf(dto.SeriesNavData.SeriesId)
	if seriesIdType != nil {
		switch seriesIdType.Kind() {
		case reflect.String:
			s := reflect.ValueOf(dto.SeriesNavData.SeriesId).String()
			if parsed, err := strconv.ParseUint(s, 10, 64); err == nil {
				seriesId = &parsed
			}
		case reflect.Float64:
			f := reflect.ValueOf(dto.SeriesNavData.SeriesId).Float()
			seriesId = utils.ToPtr(uint64(f))
		}
	}

	return &work.Work{
		Id:           utils.ParseUint64Ptr(dto.Id),
		Title:        dto.Title,
		Kind:         kind,
		Description:  parseDescription(dto.Description),
		UserId:       utils.ParseUint64Ptr(dto.UserId),
		UserName:     dto.UserName,
		Restriction:  utils.MapPtr(dto.XRestrict, work.RestrictionFromUint),
		Ai:           work.AiFromUint(utils.FromPtr(dto.AiType, work.AiDefaultUint)),
		Original:     dto.IsOriginal,
		NumPages:     dto.PageCount,
		NumViews:     dto.ViewCount,
		NumBookmarks: dto.BookmarkCount,
		NumLikes:     dto.LikeCount,
		NumComments:  dto.CommentCount,
		UploadTime:   utils.ParseLocalTimePtr(dto.CreateDate),
		DownloadTime: utils.ToPtr(downloadTime.Local()),
		SeriesId:     seriesId,
		SeriesTitle:  dto.SeriesNavData.Title,
		SeriesOrder:  dto.SeriesNavData.Order,
		Tags:         tags,
	}
}

// This function converts the work description with HTML to plain text. It does the following:
// - replaces <br> tags with line breaks
// - replaces <a> tags with their href attribute values
// - removes all other HTML tags
// - unescapes HTML entities
func parseDescription(description *string) *string {
	if description == nil {
		return nil
	}

	tagStart := 0
	textStart := 0
	builder := strings.Builder{}

	for i, r := range *description {
		switch r {
		case '<':
			tagStart = i + 1
			builder.WriteString((*description)[textStart:i])
		case '>':
			textStart = i + 1
			tag := strings.ToLower((*description)[tagStart:i])
			if strings.HasPrefix(tag, "br ") || tag == "br/" || tag == "br" {
				builder.WriteByte('\n')
			}
		}
	}

	if textStart < len(*description) {
		builder.WriteString((*description)[textStart:])
	}

	result := html.UnescapeString(builder.String())
	return &result
}
