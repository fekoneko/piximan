package work

import (
	"strconv"
	"time"
)

type ApiDto struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	IllustType    uint8  `json:"illustType"`
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
		SeriesId *string `json:"seriesId"`
		Order    *uint64 `json:"order"`
		Title    *string `json:"title"`
	} `json:"seriesNavData"`
	Tags struct {
		Tags [](struct {
			Tag string `json:"tag"`
		}) `json:"tags"`
	} `json:"tags"`
}

func (dto *ApiDto) Work(downloadTime time.Time) *Work {
	id, _ := strconv.ParseUint(dto.Id, 10, 64)
	userId, _ := strconv.ParseUint(dto.UserId, 10, 64)
	uploadTime, _ := time.Parse(dto.UploadDate, dto.UploadDate) // FIXME: parsing is broken

	var seriesId *uint64
	if dto.SeriesNavData.SeriesId != nil {
		id, _ := strconv.ParseUint(*dto.SeriesNavData.SeriesId, 10, 64)
		seriesId = &id
	}

	tags := make([]string, len(dto.Tags.Tags))
	for i, tag := range dto.Tags.Tags {
		tags[i] = tag.Tag
	}

	return &Work{
		Id:            id,
		Title:         dto.Title,
		Kind:          KindFromUint(dto.IllustType),
		Description:   dto.Description,
		UserId:        userId,
		UserName:      dto.UserName,
		Restriction:   RestrictionFromUint(dto.XRestrict),
		AiKind:        AiKindFromUint(dto.AiType),
		IsOriginal:    dto.IsOriginal,
		PageCount:     dto.PageCount,
		ViewCount:     dto.ViewCount,
		BookmarkCount: dto.BookmarkCount,
		LikeCount:     dto.LikeCount,
		CommentCount:  dto.CommentCount,
		UploadTime:    uploadTime,
		DownloadTime:  downloadTime,
		SeriesId:      seriesId,
		SeriesTitle:   dto.SeriesNavData.Title,
		SeriesOrder:   dto.SeriesNavData.Order,
		Tags:          tags,
	}
}
