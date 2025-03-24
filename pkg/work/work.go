package work

import "time"

type Work struct {
	Id            uint64
	Title         string
	Kind          Kind
	Description   string
	UserId        uint64
	UserName      string
	Restriction   Restriction
	AiKind        AiKind
	IsOriginal    bool
	PageCount     uint64
	ViewCount     uint64
	BookmarkCount uint64
	LikeCount     uint64
	CommentCount  uint64
	UploadTime    time.Time
	DownloadTime  time.Time
	SeriesId      uint64
	SeriesTitle   string
	SeriesOrder   uint64
	Tags          []string
}
