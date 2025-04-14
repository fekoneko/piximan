package work

import "time"

type Work struct {
	Id           uint64
	Title        string
	Kind         Kind
	Description  string
	UserId       uint64
	UserName     string
	Restriction  Restriction
	AiKind       AiKind
	IsOriginal   bool
	NumPages     uint64
	NumViews     uint64
	NumBookmarks uint64
	NumLikes     uint64
	NumComments  uint64
	UploadTime   time.Time
	DownloadTime time.Time
	SeriesId     *uint64
	SeriesTitle  *string
	SeriesOrder  *uint64
	Tags         []string
}
