package work

import "time"

type Work struct {
	Id           *uint64
	Title        *string
	Kind         *Kind
	Description  *string
	UserId       *uint64
	UserName     *string
	Restriction  *Restriction
	Ai           *bool
	Original     *bool
	NumPages     *uint64
	NumViews     *uint64
	NumBookmarks *uint64
	NumLikes     *uint64
	NumComments  *uint64
	UploadTime   *time.Time
	DownloadTime *time.Time
	SeriesId     *uint64
	SeriesTitle  *string
	SeriesOrder  *uint64
	Tags         *[]string
}

// Check if all fields are filled. The function doesn't report missing ai kind or series data.
func (w *Work) Full() bool {
	return w.Id != nil &&
		w.Title != nil &&
		w.Kind != nil &&
		w.Description != nil &&
		w.UserId != nil &&
		w.UserName != nil &&
		w.Restriction != nil &&
		w.Original != nil &&
		w.NumPages != nil &&
		w.NumViews != nil &&
		w.NumBookmarks != nil &&
		w.NumLikes != nil &&
		w.NumComments != nil &&
		w.UploadTime != nil &&
		w.DownloadTime != nil &&
		w.Tags != nil
}
