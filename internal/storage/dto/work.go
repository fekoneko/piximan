package dto

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

const VERSION = uint64(1)

type Work struct {
	Version     *uint64   `yaml:"_version"`
	Id          *uint64   `yaml:"id"`
	Title       *string   `yaml:"title"`
	Kind        *string   `yaml:"kind"`
	Description *string   `yaml:"description"`
	UserId      *uint64   `yaml:"user_id"`
	UserName    *string   `yaml:"user_name"`
	Restriction *string   `yaml:"restriction"`
	Ai          *bool     `yaml:"ai,omitempty"`
	Original    *bool     `yaml:"original"`
	Pages       *uint64   `yaml:"pages"`
	Views       *uint64   `yaml:"views"`
	Bookmarks   *uint64   `yaml:"bookmarks"`
	Likes       *uint64   `yaml:"likes"`
	Comments    *uint64   `yaml:"comments"`
	Uploaded    *string   `yaml:"uploaded"`
	Downloaded  *string   `yaml:"downloaded"`
	SeriesId    *uint64   `yaml:"series_id,omitempty"`
	SeriesTitle *string   `yaml:"series_title,omitempty"`
	SeriesOrder *uint64   `yaml:"series_order,omitempty"`
	Tags        *[]string `yaml:"tags"`
}

func ToDto(w *work.Work) *Work {
	return &Work{
		Version:     utils.ToPtr(VERSION),
		Id:          w.Id,
		Title:       w.Title,
		Kind:        utils.MapPtr(w.Kind, work.Kind.String),
		Description: w.Description,
		UserId:      w.UserId,
		UserName:    w.UserName,
		Restriction: utils.MapPtr(w.Restriction, work.Restriction.String),
		Ai:          w.AiKind.Bool(),
		Original:    w.IsOriginal,
		Pages:       w.NumPages,
		Views:       w.NumViews,
		Bookmarks:   w.NumBookmarks,
		Likes:       w.NumLikes,
		Comments:    w.NumComments,
		Uploaded:    utils.FormatUTCTimePtr(w.UploadTime),
		Downloaded:  utils.FormatUTCTimePtr(w.DownloadTime),
		SeriesId:    w.SeriesId,
		SeriesTitle: w.SeriesTitle,
		SeriesOrder: w.SeriesOrder,
		Tags:        w.Tags,
	}
}
