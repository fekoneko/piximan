package dto

import (
	"fmt"
	"regexp"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/utils"
)

type Rules struct {
	Ids                    *[]uint64 `yaml:"ids"`
	NotIds                 *[]uint64 `yaml:"not_ids"`
	TitleContains          *[]string `yaml:"title_contains"`
	TitleNotContains       *[]string `yaml:"title_not_contains"`
	TitleRegexp            *string   `yaml:"title_regexp"`
	Kinds                  *[]string `yaml:"kinds"`
	DescriptionContains    *[]string `yaml:"description_contains"`
	DescriptionNotContains *[]string `yaml:"description_not_contains"`
	DescriptionRegexp      *string   `yaml:"description_regexp"`
	UserIds                *[]uint64 `yaml:"user_ids"`
	NotUserIds             *[]uint64 `yaml:"not_user_ids"`
	UserNames              *[]string `yaml:"user_names"`
	NotUserNames           *[]string `yaml:"not_user_names"`
	Restrictions           *[]string `yaml:"restrictions"`
	Ai                     *bool     `yaml:"ai"`
	Original               *bool     `yaml:"original"`
	PagesLessThan          *uint64   `yaml:"pages_less_than"`
	PagesMoreThan          *uint64   `yaml:"pages_more_than"`
	ViewsLessThan          *uint64   `yaml:"views_less_than"`
	ViewsMoreThan          *uint64   `yaml:"views_more_than"`
	BookmarksLessThan      *uint64   `yaml:"bookmarks_less_than"`
	BookmarksMoreThan      *uint64   `yaml:"bookmarks_more_than"`
	LikesLessThan          *uint64   `yaml:"likes_less_than"`
	LikesMoreThan          *uint64   `yaml:"likes_more_than"`
	CommentsLessThan       *uint64   `yaml:"comments_less_than"`
	CommentsMoreThan       *uint64   `yaml:"comments_more_than"`
	UploadedBefore         *string   `yaml:"uploaded_before"`
	UploadedAfter          *string   `yaml:"uploaded_after"`
	Series                 *bool     `yaml:"series"`
	SeriesIds              *[]uint64 `yaml:"series_ids"`
	NotSeriesIds           *[]uint64 `yaml:"not_series_ids"`
	SeriesTitleContains    *[]string `yaml:"series_title_contains"`
	SeriesTitleNotContains *[]string `yaml:"series_title_not_contains"`
	SeriesTitleRegexp      *string   `yaml:"series_title_regexp"`
	Tags                   *[]string `yaml:"tags"`
	NotTags                *[]string `yaml:"not_tags"`
}

func (dto *Rules) FromDto() (*queue.Rules, error) {
	var titleRegexp *regexp.Regexp
	if dto.TitleRegexp != nil {
		var err error
		if titleRegexp, err = regexp.Compile(*dto.TitleRegexp); err != nil {
			return nil, fmt.Errorf("incorrect value for rule 'title_regexp': %v", err)
		}
	}

	var kinds *[]work.Kind
	if dto.Kinds != nil {
		kinds = utils.ToPtr(make([]work.Kind, len(*dto.Kinds)))
		for i, kind := range *dto.Kinds {
			if !work.ValidKindString(kind) {
				return nil, fmt.Errorf("incorrect value for rule 'kinds': invalid work kind: %v", kind)
			}
			(*kinds)[i] = work.KindFromString(kind)
		}
	}

	var descriptionRegexp *regexp.Regexp
	if dto.DescriptionRegexp != nil {
		var err error
		if descriptionRegexp, err = regexp.Compile(*dto.DescriptionRegexp); err != nil {
			return nil, fmt.Errorf("incorrect value for rule 'description_regexp': %v", err)
		}
	}

	var restrictions *[]work.Restriction
	if dto.Restrictions != nil {
		restrictions = utils.ToPtr(make([]work.Restriction, len(*dto.Restrictions)))
		for i, restriction := range *dto.Restrictions {
			if !work.ValidRestrictionString(restriction) {
				return nil, fmt.Errorf("incorrect value for rule 'restrictions': invalid restriction: %v", restriction)
			}
			(*restrictions)[i] = work.RestrictionFromString(restriction)
		}
	}

	uploadedBefore, err := utils.ParseLocalTimePtrStrict(dto.UploadedBefore)
	if err != nil {
		return nil, fmt.Errorf("incorrect value for rule 'uploaded_before': %v", err)
	}

	uploadedAfter, err := utils.ParseLocalTimePtrStrict(dto.UploadedAfter)
	if err != nil {
		return nil, fmt.Errorf("incorrect value for rule 'uploaded_after': %v", err)
	}

	var seriesTitleRegexp *regexp.Regexp
	if dto.SeriesTitleRegexp != nil {
		var err error
		if seriesTitleRegexp, err = regexp.Compile(*dto.SeriesTitleRegexp); err != nil {
			return nil, fmt.Errorf("incorrect value for rule 'series_title_regexp': %v", err)
		}
	}

	rules := &queue.Rules{
		Ids:                    dto.Ids,
		NotIds:                 dto.NotIds,
		TitleContains:          dto.TitleContains,
		TitleNotContains:       dto.TitleNotContains,
		TitleRegexp:            titleRegexp,
		Kinds:                  kinds,
		DescriptionContains:    dto.DescriptionContains,
		DescriptionNotContains: dto.DescriptionNotContains,
		DescriptionRegexp:      descriptionRegexp,
		UserIds:                dto.UserIds,
		NotUserIds:             dto.NotUserIds,
		UserNames:              dto.UserNames,
		NotUserNames:           dto.NotUserNames,
		Restrictions:           restrictions,
		Ai:                     dto.Ai,
		Original:               dto.Original,
		PagesLessThan:          dto.PagesLessThan,
		PagesMoreThan:          dto.PagesMoreThan,
		ViewsLessThan:          dto.ViewsLessThan,
		ViewsMoreThan:          dto.ViewsMoreThan,
		BookmarksLessThan:      dto.BookmarksLessThan,
		BookmarksMoreThan:      dto.BookmarksMoreThan,
		LikesLessThan:          dto.LikesLessThan,
		LikesMoreThan:          dto.LikesMoreThan,
		CommentsLessThan:       dto.CommentsLessThan,
		CommentsMoreThan:       dto.CommentsMoreThan,
		UploadedBefore:         uploadedBefore,
		UploadedAfter:          uploadedAfter,
		Series:                 dto.Series,
		SeriesIds:              dto.SeriesIds,
		NotSeriesIds:           dto.NotSeriesIds,
		SeriesTitleContains:    dto.SeriesTitleContains,
		SeriesTitleNotContains: dto.SeriesTitleNotContains,
		SeriesTitleRegexp:      seriesTitleRegexp,
		Tags:                   dto.Tags,
		NotTags:                dto.NotTags,
	}

	return rules, nil
}
