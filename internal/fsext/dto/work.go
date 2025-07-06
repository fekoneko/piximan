package dto

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

const VERSION = uint64(1)

type Work struct {
	Version     *uint64   `yaml:"_version,omitempty"`
	Id          *uint64   `yaml:"id,omitempty"`
	Title       *string   `yaml:"title,omitempty"`
	Kind        *string   `yaml:"kind,omitempty"`
	Description *string   `yaml:"description,omitempty"`
	UserId      *uint64   `yaml:"user_id,omitempty"`
	UserName    *string   `yaml:"user_name,omitempty"`
	Restriction *string   `yaml:"restriction,omitempty"`
	Ai          *bool     `yaml:"ai,omitempty"`
	Original    *bool     `yaml:"original,omitempty"`
	Pages       *uint64   `yaml:"pages,omitempty"`
	Views       *uint64   `yaml:"views,omitempty"`
	Bookmarks   *uint64   `yaml:"bookmarks,omitempty"`
	Likes       *uint64   `yaml:"likes,omitempty"`
	Comments    *uint64   `yaml:"comments,omitempty"`
	Uploaded    *string   `yaml:"uploaded,omitempty"`
	Downloaded  *string   `yaml:"downloaded,omitempty"`
	SeriesId    *uint64   `yaml:"series_id,omitempty"`
	SeriesTitle *string   `yaml:"series_title,omitempty"`
	SeriesOrder *uint64   `yaml:"series_order,omitempty"`
	Tags        *[]string `yaml:"tags,omitempty"`
}

func WorkToDto(w *work.Work) *Work {
	return &Work{
		Version:     utils.ToPtr(VERSION),
		Id:          w.Id,
		Title:       w.Title,
		Kind:        utils.MapPtr(w.Kind, work.Kind.String),
		Description: w.Description,
		UserId:      w.UserId,
		UserName:    w.UserName,
		Restriction: utils.MapPtr(w.Restriction, work.Restriction.String),
		Ai:          w.Ai,
		Original:    w.Original,
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

func (dto *Work) FromDto() (w *work.Work, warning error) {
	if dto.Version == nil {
		warning = fmt.Errorf("metadata version is missing")
	} else if *dto.Version != VERSION {
		warning = fmt.Errorf("metadata version mismatch: expected %v, got %v", VERSION, *dto.Version)
	}

	// TODO: warn if some fields are incorrect or there are extra fields

	w = &work.Work{
		Id:           dto.Id,
		Title:        dto.Title,
		Kind:         utils.MapPtr(dto.Kind, work.KindFromString),
		Description:  dto.Description,
		UserId:       dto.UserId,
		UserName:     dto.UserName,
		Restriction:  utils.MapPtr(dto.Restriction, work.RestrictionFromString),
		Ai:           dto.Ai,
		Original:     dto.Original,
		NumPages:     dto.Pages,
		NumViews:     dto.Views,
		NumBookmarks: dto.Bookmarks,
		NumLikes:     dto.Likes,
		NumComments:  dto.Comments,
		UploadTime:   utils.ParseLocalTimePtr(dto.Uploaded),
		DownloadTime: utils.ParseLocalTimePtr(dto.Downloaded),
		SeriesId:     dto.SeriesId,
		SeriesTitle:  dto.SeriesTitle,
		SeriesOrder:  dto.SeriesOrder,
		Tags:         dto.Tags,
	}

	return w, warning
}
