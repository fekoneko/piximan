package dto

import (
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

type BookmarkWork struct {
	Id          *string   `json:"id"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	UserId      *string   `json:"userId"`
	UserName    *string   `json:"userName"`
	XRestrict   *uint8    `json:"xRestrict"`
	AiType      *uint8    `json:"aiType"`
	PageCount   *uint64   `json:"pageCount"`
	CreateDate  *string   `json:"createDate"`
	Tags        *[]string `json:"tags"`
}

func (dto *BookmarkWork) FromDto(kind *work.Kind, downloadTime time.Time) (*work.Work, *time.Time) {
	bookmarkedTime := utils.ParseLocalTimePtr(dto.CreateDate)

	work := &work.Work{
		Id:           utils.ParseUint64Ptr(dto.Id),
		Title:        dto.Title,
		Kind:         kind,
		Description:  formatDescription(dto.Description),
		UserId:       utils.ParseUint64Ptr(dto.UserId),
		UserName:     dto.UserName,
		Restriction:  utils.MapPtr(dto.XRestrict, work.RestrictionFromUint),
		AiKind:       utils.MapPtr(dto.XRestrict, work.AiKindFromUint),
		NumPages:     dto.PageCount,
		DownloadTime: utils.ToPtr(downloadTime.Local()),
		Tags:         dto.Tags,
	}

	return work, bookmarkedTime
}
