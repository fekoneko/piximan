package dto

import (
	"strconv"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

type BookmarkWork struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	UserId      string   `json:"userId"`
	UserName    string   `json:"userName"`
	XRestrict   uint8    `json:"xRestrict"`
	AiType      uint8    `json:"aiType"`
	PageCount   uint64   `json:"pageCount"`
	CreateDate  string   `json:"createDate"`
	Tags        []string `json:"tags"`
}

func (dto *BookmarkWork) FromDto(kind work.Kind, downloadTime time.Time) (*work.Work, *time.Time) {
	id, _ := strconv.ParseUint(dto.Id, 10, 64)
	userId, _ := strconv.ParseUint(dto.UserId, 10, 64)

	bookmarkedTime, err := time.Parse(time.RFC3339, dto.CreateDate)
	localBookmarkedTime := bookmarkedTime.Local()
	localBookmarkedTimePtr := &localBookmarkedTime
	if err != nil {
		localBookmarkedTimePtr = nil
	}

	work := &work.Work{
		Id:           id,
		Title:        dto.Title,
		Kind:         kind,
		Description:  formatDescription(dto.Description),
		UserId:       userId,
		UserName:     dto.UserName,
		Restriction:  work.RestrictionFromUint(dto.XRestrict),
		AiKind:       work.AiKindFromUint(dto.AiType),
		NumPages:     dto.PageCount,
		DownloadTime: downloadTime.Local(),
		Tags:         dto.Tags,
	}

	return work, localBookmarkedTimePtr
}
