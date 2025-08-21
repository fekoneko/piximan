package dto

import (
	"fmt"
	"regexp"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/utils"
)

const RulesVersion = uint64(1)

type Rules struct {
	Version                *uint64   `yaml:"_version,omitempty"`
	Ids                    *[]uint64 `yaml:"ids,omitempty"`
	NotIds                 *[]uint64 `yaml:"not_ids,omitempty"`
	TitleContains          *[]string `yaml:"title_contains,omitempty"`
	TitleNotContains       *[]string `yaml:"title_not_contains,omitempty"`
	TitleRegexp            *string   `yaml:"title_regexp,omitempty"`
	Kinds                  *[]string `yaml:"kinds,omitempty"`
	DescriptionContains    *[]string `yaml:"description_contains,omitempty"`
	DescriptionNotContains *[]string `yaml:"description_not_contains,omitempty"`
	DescriptionRegexp      *string   `yaml:"description_regexp,omitempty"`
	UserIds                *[]uint64 `yaml:"user_ids,omitempty"`
	NotUserIds             *[]uint64 `yaml:"not_user_ids,omitempty"`
	UserNames              *[]string `yaml:"user_names,omitempty"`
	NotUserNames           *[]string `yaml:"not_user_names,omitempty"`
	Restrictions           *[]string `yaml:"restrictions,omitempty"`
	Ai                     *bool     `yaml:"ai,omitempty"`
	Original               *bool     `yaml:"original,omitempty"`
	PagesLessThan          *uint64   `yaml:"pages_less_than,omitempty"`
	PagesMoreThan          *uint64   `yaml:"pages_more_than,omitempty"`
	ViewsLessThan          *uint64   `yaml:"views_less_than,omitempty"`
	ViewsMoreThan          *uint64   `yaml:"views_more_than,omitempty"`
	BookmarksLessThan      *uint64   `yaml:"bookmarks_less_than,omitempty"`
	BookmarksMoreThan      *uint64   `yaml:"bookmarks_more_than,omitempty"`
	LikesLessThan          *uint64   `yaml:"likes_less_than,omitempty"`
	LikesMoreThan          *uint64   `yaml:"likes_more_than,omitempty"`
	CommentsLessThan       *uint64   `yaml:"comments_less_than,omitempty"`
	CommentsMoreThan       *uint64   `yaml:"comments_more_than,omitempty"`
	UploadedBefore         *string   `yaml:"uploaded_before,omitempty"`
	UploadedAfter          *string   `yaml:"uploaded_after,omitempty"`
	Series                 *bool     `yaml:"series,omitempty"`
	SeriesIds              *[]uint64 `yaml:"series_ids,omitempty"`
	NotSeriesIds           *[]uint64 `yaml:"not_series_ids,omitempty"`
	SeriesTitleContains    *[]string `yaml:"series_title_contains,omitempty"`
	SeriesTitleNotContains *[]string `yaml:"series_title_not_contains,omitempty"`
	SeriesTitleRegexp      *string   `yaml:"series_title_regexp,omitempty"`
	Tags                   *[]string `yaml:"tags,omitempty"`
	NotTags                *[]string `yaml:"not_tags,omitempty"`
}

func RulesToDto(r *rules.Rules) *Rules {
	return &Rules{} // TODO: implement
}

// Missing version field is not considered warning as we're not forcing user to specify it.
func (dto *Rules) FromDto() (r *rules.Rules, warning error, err error) {
	if dto.Version != nil && *dto.Version != WorkVersion {
		warning = fmt.Errorf("download rules version mismatch: expected %v, got %v", WorkVersion, *dto.Version)
	}

	var titleRegexp *regexp.Regexp
	if dto.TitleRegexp != nil {
		var err error
		if titleRegexp, err = regexp.Compile(*dto.TitleRegexp); err != nil {
			return nil, warning, fmt.Errorf("incorrect value for rule 'title_regexp': %v", err)
		}
	}

	var kinds *[]work.Kind
	if dto.Kinds != nil {
		kinds = utils.ToPtr(make([]work.Kind, len(*dto.Kinds)))
		for i, kind := range *dto.Kinds {
			if !work.ValidKindString(kind) {
				return nil, warning, fmt.Errorf("incorrect value for rule 'kinds': invalid work kind: %v", kind)
			}
			(*kinds)[i] = work.KindFromString(kind)
		}
	}

	var descriptionRegexp *regexp.Regexp
	if dto.DescriptionRegexp != nil {
		var err error
		if descriptionRegexp, err = regexp.Compile(*dto.DescriptionRegexp); err != nil {
			return nil, warning, fmt.Errorf("incorrect value for rule 'description_regexp': %v", err)
		}
	}

	var restrictions *[]work.Restriction
	if dto.Restrictions != nil {
		restrictions = utils.ToPtr(make([]work.Restriction, len(*dto.Restrictions)))
		for i, restriction := range *dto.Restrictions {
			if !work.ValidRestrictionString(restriction) {
				return nil, warning, fmt.Errorf("incorrect value for rule 'restrictions': invalid restriction: %v", restriction)
			}
			(*restrictions)[i] = work.RestrictionFromString(restriction)
		}
	}

	uploadedBefore, err := utils.ParseLocalTimePtrStrict(dto.UploadedBefore)
	if err != nil {
		return nil, warning, fmt.Errorf("incorrect value for rule 'uploaded_before': %v", err)
	}

	uploadedAfter, err := utils.ParseLocalTimePtrStrict(dto.UploadedAfter)
	if err != nil {
		return nil, warning, fmt.Errorf("incorrect value for rule 'uploaded_after': %v", err)
	}

	var seriesTitleRegexp *regexp.Regexp
	if dto.SeriesTitleRegexp != nil {
		var err error
		if seriesTitleRegexp, err = regexp.Compile(*dto.SeriesTitleRegexp); err != nil {
			return nil, warning, fmt.Errorf("incorrect value for rule 'series_title_regexp': %v", err)
		}
	}

	rules := &rules.Rules{
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

	return rules, warning, nil
}
