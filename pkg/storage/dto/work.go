package dto

import (
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

const version = 1

type Work struct {
	Version     uint64   `yaml:"_version"`
	Id          uint64   `yaml:"id"`
	Title       string   `yaml:"title"`
	Kind        string   `yaml:"kind"`
	Description string   `yaml:"description"`
	UserId      uint64   `yaml:"user_id"`
	UserName    string   `yaml:"user_name"`
	Restriction string   `yaml:"restriction"`
	Ai          *bool    `yaml:"ai,omitempty"`
	Original    bool     `yaml:"original"`
	Pages       uint64   `yaml:"pages"`
	Views       uint64   `yaml:"views"`
	Bookmarks   uint64   `yaml:"bookmarks"`
	Likes       uint64   `yaml:"likes"`
	Comments    uint64   `yaml:"comments"`
	Uploaded    string   `yaml:"uploaded"`
	Downloaded  string   `yaml:"downloaded"`
	SeriesId    *uint64  `yaml:"series_id,omitempty"`
	SeriesTitle *string  `yaml:"series_title,omitempty"`
	SeriesOrder *uint64  `yaml:"series_order,omitempty"`
	Tags        []string `yaml:"tags"`
}

func ToDto(work *work.Work) *Work {
	return &Work{
		Version:     version,
		Id:          work.Id,
		Title:       work.Title,
		Kind:        work.Kind.String(),
		Description: work.Description,
		UserId:      work.UserId,
		UserName:    work.UserName,
		Restriction: work.Restriction.String(),
		Ai:          work.AiKind.Bool(),
		Original:    work.IsOriginal,
		Pages:       work.PageCount,
		Views:       work.ViewCount,
		Bookmarks:   work.BookmarkCount,
		Likes:       work.LikeCount,
		Comments:    work.CommentCount,
		Uploaded:    work.UploadTime.UTC().Format(time.RFC3339),
		Downloaded:  work.DownloadTime.UTC().Format(time.RFC3339),
		SeriesId:    work.SeriesId,
		SeriesTitle: work.SeriesTitle,
		SeriesOrder: work.SeriesOrder,
		Tags:        work.Tags,
	}
}
