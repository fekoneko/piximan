package dto

import (
	"reflect"
	"strconv"
	"time"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
)

type BookmarkWork struct {
	Id          any       `json:"id"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	UserId      any       `json:"userId"`
	UserName    *string   `json:"userName"`
	XRestrict   *uint8    `json:"xRestrict"`
	AiType      *uint8    `json:"aiType"`
	PageCount   *uint64   `json:"pageCount"`
	CreateDate  *string   `json:"createDate"`
	Tags        *[]string `json:"tags"`
}

func (dto *BookmarkWork) FromDto(
	kind *work.Kind, downloadTime time.Time,
) (w *work.Work, unlisted bool) {
	var id *uint64
	idTypeKind := reflect.TypeOf(dto.Id).Kind()
	switch idTypeKind {
	case reflect.String:
		if parsed, err := strconv.ParseUint(reflect.ValueOf(dto.Id).String(), 10, 64); err == nil {
			id = &parsed
		}
	case reflect.Float64:
		parsed := reflect.ValueOf(dto.Id).Float()
		id := utils.ToPtr(uint64(parsed))
		work := work.Work{Id: id}
		return &work, true
	}

	var userId *uint64
	userIdTypeKind := reflect.TypeOf(dto.UserId).Kind()
	if userIdTypeKind == reflect.String {
		if parsed, err := strconv.ParseUint(reflect.ValueOf(dto.UserId).String(), 10, 64); err == nil {
			userId = &parsed
		}
	}

	work := &work.Work{
		Id:           id,
		Title:        dto.Title,
		Kind:         kind,
		Description:  formatDescription(dto.Description),
		UserId:       userId,
		UserName:     dto.UserName,
		Restriction:  utils.MapPtr(dto.XRestrict, work.RestrictionFromUint),
		AiKind:       utils.MapPtr(dto.XRestrict, work.AiKindFromUint),
		NumPages:     dto.PageCount,
		UploadTime:   utils.ParseLocalTimePtr(dto.CreateDate),
		DownloadTime: utils.ToPtr(downloadTime.Local()),
		Tags:         dto.Tags,
	}

	return work, false
}
