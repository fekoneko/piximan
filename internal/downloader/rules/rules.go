package rules

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
)

type Rules struct {
	Ids                    *[]uint64           `yaml:"ids"`
	NotIds                 *[]uint64           `yaml:"not_ids"`
	TitleContains          *[]string           `yaml:"title_contains"`
	TitleNotContains       *[]string           `yaml:"title_not_contains"`
	TitleRegexp            *regexp.Regexp      `yaml:"title_regexp"`
	Kinds                  *[]work.Kind        `yaml:"kinds"`
	DescriptionContains    *[]string           `yaml:"description_contains"`
	DescriptionNotContains *[]string           `yaml:"description_not_contains"`
	DescriptionRegexp      *regexp.Regexp      `yaml:"description_regexp"`
	UserIds                *[]uint64           `yaml:"user_ids"`
	NotUserIds             *[]uint64           `yaml:"not_user_ids"`
	UserNames              *[]string           `yaml:"user_names"`
	NotUserNames           *[]string           `yaml:"not_user_names"`
	Restrictions           *[]work.Restriction `yaml:"restrictions"`
	Ai                     *bool               `yaml:"ai"`
	Original               *bool               `yaml:"original"`
	PagesLessThan          *uint64             `yaml:"pages_less_than"`
	PagesMoreThan          *uint64             `yaml:"pages_more_than"`
	ViewsLessThan          *uint64             `yaml:"views_less_than"`
	ViewsMoreThan          *uint64             `yaml:"views_more_than"`
	BookmarksLessThan      *uint64             `yaml:"bookmarks_less_than"`
	BookmarksMoreThan      *uint64             `yaml:"bookmarks_more_than"`
	LikesLessThan          *uint64             `yaml:"likes_less_than"`
	LikesMoreThan          *uint64             `yaml:"likes_more_than"`
	CommentsLessThan       *uint64             `yaml:"comments_less_than"`
	CommentsMoreThan       *uint64             `yaml:"comments_more_than"`
	UploadedBefore         *time.Time          `yaml:"uploaded_before"`
	UploadedAfter          *time.Time          `yaml:"uploaded_after"`
	Series                 *bool               `yaml:"series"`
	SeriesIds              *[]uint64           `yaml:"series_ids"`
	NotSeriesIds           *[]uint64           `yaml:"not_series_ids"`
	SeriesTitleContains    *[]string           `yaml:"series_title_contains"`
	SeriesTitleNotContains *[]string           `yaml:"series_title_not_contains"`
	SeriesTitleRegexp      *regexp.Regexp      `yaml:"series_title_regexp"`
	Tags                   *[]string           `yaml:"tags"`
	NotTags                *[]string           `yaml:"not_tags"`
}

// Checks weather the artwork is worth getting metadata for a full MatchWork() call.
func (r *Rules) MatchArtworkId(id uint64) bool {
	return (r.Ids == nil || slices.Contains(*r.Ids, id)) &&
		(r.NotIds == nil || !slices.Contains(*r.NotIds, id)) &&
		(r.Kinds == nil || slices.ContainsFunc(*r.Kinds, func(kind work.Kind) bool {
			return kind == work.KindIllust || kind == work.KindManga || kind == work.KindUgoira
		}))
}

// Checks weather the novel is worth getting metadata for a full MatchWork() call.
func (r *Rules) MatchNovelId(id uint64) bool {
	return (r.Ids == nil || slices.Contains(*r.Ids, id)) &&
		(r.NotIds == nil || !slices.Contains(*r.NotIds, id)) &&
		(r.Kinds == nil || slices.Contains(*r.Kinds, work.KindNovel))
}

// Checkes weather the work matches the rules and can be downloaded.
// Pass partial if the work series data is unknown (e.g. the work was received from bookmarks request).
// In this case additional warning will be returned if some series rule is defined.
func (r *Rules) MatchWork(w *work.Work, partial bool) (matches bool, warnings []error) {
	matches = matchManyToOne(r.Ids, w.Id, "ids", "id", &warnings) &&
		matchManyToOneNot(r.NotIds, w.Id, "not_ids", "id", &warnings) &&
		matchSubstrings(r.TitleContains, w.Title, "title_contains", "title", &warnings) &&
		matchSubstringsNot(r.TitleNotContains, w.Title, "title_not_contains", "title", &warnings) &&
		matchRegexp(r.TitleRegexp, w.Title, "title_regexp", "title", &warnings) &&
		matchManyToOne(r.Kinds, w.Kind, "kinds", "kind", &warnings) &&
		matchSubstrings(r.DescriptionContains, w.Description, "description_contains", "description", &warnings) &&
		matchSubstringsNot(r.DescriptionNotContains, w.Description, "description_not_contains", "description", &warnings) &&
		matchRegexp(r.DescriptionRegexp, w.Description, "description_regexp", "description", &warnings) &&
		matchManyToOne(r.UserIds, w.UserId, "user_ids", "user id", &warnings) &&
		matchManyToOneNot(r.NotUserIds, w.UserId, "not_user_ids", "user id", &warnings) &&
		matchManyToOne(r.UserNames, w.UserName, "user_names", "user name", &warnings) &&
		matchManyToOneNot(r.NotUserNames, w.UserName, "not_user_names", "user name", &warnings) &&
		matchManyToOne(r.Restrictions, w.Restriction, "restrictions", "restriction", &warnings) &&
		matchOneToOne(r.Ai, w.Ai, "ai", "ai kind", &warnings) &&
		matchOneToOne(r.Original, w.Original, "original", "original", &warnings) &&
		matchLessThan(r.PagesLessThan, w.NumPages, "pages_less_than", "pages count", &warnings) &&
		matchMoreThan(r.PagesMoreThan, w.NumPages, "pages_more_than", "pages count", &warnings) &&
		matchLessThan(r.ViewsLessThan, w.NumViews, "views_less_than", "views count", &warnings) &&
		matchMoreThan(r.ViewsMoreThan, w.NumViews, "views_more_than", "views count", &warnings) &&
		matchLessThan(r.BookmarksLessThan, w.NumBookmarks, "bookmarks_less_than", "bookmarks count", &warnings) &&
		matchMoreThan(r.BookmarksMoreThan, w.NumBookmarks, "bookmarks_more_than", "bookmarks count", &warnings) &&
		matchLessThan(r.LikesLessThan, w.NumLikes, "likes_less_than", "likes count", &warnings) &&
		matchMoreThan(r.LikesMoreThan, w.NumLikes, "likes_more_than", "likes count", &warnings) &&
		matchLessThan(r.CommentsLessThan, w.NumComments, "comments_less_than", "comments count", &warnings) &&
		matchMoreThan(r.CommentsMoreThan, w.NumComments, "comments_more_than", "comments count", &warnings) &&
		matchBefore(r.UploadedBefore, w.UploadTime, "uploaded_before", "upload time", &warnings) &&
		matchAfter(r.UploadedAfter, w.UploadTime, "uploaded_after", "upload time", &warnings) &&
		matchDefined(r.Series, w.SeriesId, "series", "series data", &warnings, partial) &&
		matchManyToOne(r.SeriesIds, w.SeriesId, "series_ids", "series id", &warnings) &&
		matchManyToOneNot(r.NotSeriesIds, w.SeriesId, "not_series_ids", "series id", &warnings) &&
		matchSubstrings(r.SeriesTitleContains, w.SeriesTitle, "series_title_contains", "series title", &warnings) &&
		matchSubstringsNot(r.SeriesTitleNotContains, w.SeriesTitle, "series_title_not_contains", "series title", &warnings) &&
		matchRegexp(r.SeriesTitleRegexp, w.SeriesTitle, "series_title_regexp", "series title", &warnings) &&
		matchManyToMany(r.Tags, w.Tags, "tags", "tags", &warnings) &&
		matchManyToManyNot(r.NotTags, w.Tags, "not_tags", "tags", &warnings)

	if !matches {
		warnings = []error{}
	}
	return matches, warnings
}

func matchOneToOne[T comparable](r *T, f *T, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || *r == *f
}

func matchManyToOne[T comparable](r *[]T, f *T, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || slices.Contains(*r, *f)
}

func matchManyToOneNot[T comparable](r *[]T, f *T, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || !slices.Contains(*r, *f)
}

func matchManyToMany[T comparable](r *[]T, f *[]T, rule string, field string, warnings *[]error) bool {
	if skipped(r, f, rule, field, warnings, true) {
		return true
	}
	for _, value := range *r {
		if slices.Contains(*f, value) {
			return true
		}
	}
	return false
}

func matchManyToManyNot[T comparable](r *[]T, f *[]T, rule string, field string, warnings *[]error) bool {
	if skipped(r, f, rule, field, warnings, true) {
		return true
	}
	for _, value := range *r {
		if slices.Contains(*f, value) {
			return false
		}
	}
	return true
}

func matchSubstrings(r *[]string, f *string, rule string, field string, warnings *[]error) bool {
	if skipped(r, f, rule, field, warnings, true) {
		return true
	}
	for _, value := range *r {
		if strings.Contains(strings.ToLower(*f), strings.ToLower(value)) {
			return true
		}
	}
	return false
}

func matchSubstringsNot(r *[]string, f *string, rule string, field string, warnings *[]error) bool {
	if skipped(r, f, rule, field, warnings, true) {
		return true
	}
	for _, value := range *r {
		if strings.Contains(strings.ToLower(*f), strings.ToLower(value)) {
			return false
		}
	}
	return true
}

func matchRegexp(r *regexp.Regexp, f *string, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || r.MatchString(*f)
}

func matchLessThan(r *uint64, f *uint64, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || *f < *r
}

func matchMoreThan(r *uint64, f *uint64, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || *f > *r
}

func matchBefore(r *time.Time, f *time.Time, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || f.Before(*r)
}

func matchAfter(r *time.Time, f *time.Time, rule string, field string, warnings *[]error) bool {
	skipped := skipped(r, f, rule, field, warnings, true)
	return skipped || f.After(*r)
}

func matchDefined[T comparable](
	r *bool, f *T, rule string, field string, warnings *[]error, skipUndefinedField bool,
) bool {
	skipped := skipped(r, f, rule, field, warnings, skipUndefinedField)
	return skipped || (*r && f != nil) || (!*r && f == nil)
}

func skipped[R any, F any](r *R, f *F, rule string, field string, warnings *[]error, skipUndefinedField bool) bool {
	if r == nil {
		return true
	} else if f == nil && skipUndefinedField {
		*warnings = append(*warnings, fmt.Errorf("ignoring rule '%v': %v is unknown", rule, field))
		return true
	}
	return false
}

func (r *Rules) Count() int {
	count := 0

	countRule(r.Ids, &count)
	countRule(r.NotIds, &count)
	countRule(r.TitleContains, &count)
	countRule(r.TitleNotContains, &count)
	countRule(r.TitleRegexp, &count)
	countRule(r.Kinds, &count)
	countRule(r.DescriptionContains, &count)
	countRule(r.DescriptionNotContains, &count)
	countRule(r.DescriptionRegexp, &count)
	countRule(r.UserIds, &count)
	countRule(r.NotUserIds, &count)
	countRule(r.UserNames, &count)
	countRule(r.NotUserNames, &count)
	countRule(r.Restrictions, &count)
	countRule(r.Ai, &count)
	countRule(r.Original, &count)
	countRule(r.PagesLessThan, &count)
	countRule(r.PagesMoreThan, &count)
	countRule(r.ViewsLessThan, &count)
	countRule(r.ViewsMoreThan, &count)
	countRule(r.BookmarksLessThan, &count)
	countRule(r.BookmarksMoreThan, &count)
	countRule(r.LikesLessThan, &count)
	countRule(r.LikesMoreThan, &count)
	countRule(r.CommentsLessThan, &count)
	countRule(r.CommentsMoreThan, &count)
	countRule(r.UploadedBefore, &count)
	countRule(r.UploadedAfter, &count)
	countRule(r.Series, &count)
	countRule(r.SeriesIds, &count)
	countRule(r.NotSeriesIds, &count)
	countRule(r.SeriesTitleContains, &count)
	countRule(r.SeriesTitleNotContains, &count)
	countRule(r.SeriesTitleRegexp, &count)
	countRule(r.Tags, &count)
	countRule(r.NotTags, &count)

	return count
}

func countRule[T any](r *T, count *int) {
	if r != nil {
		*count++
	}
}
