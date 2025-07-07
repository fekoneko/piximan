package dto

import (
	"reflect"
	"strconv"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
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
	idType := reflect.TypeOf(dto.Id)
	if idType != nil {
		switch idType.Kind() {
		case reflect.String:
			s := reflect.ValueOf(dto.Id).String()
			if parsed, err := strconv.ParseUint(s, 10, 64); err == nil {
				id = &parsed
			}
		case reflect.Float64:
			f := reflect.ValueOf(dto.Id).Float()
			work := work.Work{Id: utils.ToPtr(uint64(f))}
			return &work, true
		}
	}

	var userId *uint64
	userIdType := reflect.TypeOf(dto.UserId)
	if userIdType != nil && userIdType.Kind() == reflect.String {
		s := reflect.ValueOf(dto.UserId).String()
		if parsed, err := strconv.ParseUint(s, 10, 64); err == nil {
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
		Ai:           work.AiFromUint(utils.FromPtr(dto.AiType, work.AiDefaultUint)),
		NumPages:     dto.PageCount,
		UploadTime:   utils.ParseLocalTimePtr(dto.CreateDate),
		DownloadTime: utils.ToPtr(downloadTime.Local()),
		Tags:         dto.Tags,
	}

	return work, false
}
