package fsext

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

func FormatWorkPath(pattern string, w *work.Work) (string, error) {
	path, err := filepath.Abs(pattern)
	if err != nil {
		return "", err
	}
	sections := strings.Split(path, string(filepath.Separator))
	if len(sections) != 0 && len(sections[0]) == 0 {
		sections[0] = string(filepath.Separator)
	}

	replacer := workPathReplacer(w)
	for i, section := range sections {
		if i == 0 || section == "." || section == ".." {
			continue
		}
		filename := replacer.Replace(section)
		sections[i] = FormatFilename(filename)
	}

	return filepath.Join(sections...), nil
}

func FormatWorkPaths(patterns []string, w *work.Work) ([]string, error) {
	paths := make([]string, len(patterns))
	for i, pattern := range patterns {
		path, err := FormatWorkPath(pattern, w)
		if err != nil {
			return nil, err
		}
		paths[i] = path
	}
	return paths, nil
}

func WorkPathValid(pattern string) error {
	for _, match := range pathSubstitutionRegexp.FindAllString(pattern, -1) {
		if _, ok := workPathSubstitutions[match]; !ok {
			return fmt.Errorf("pattern contains unknown substitution %q", match)
		}
	}
	return nil
}

func workPathReplacer(w *work.Work) *strings.Replacer {
	oldNew := make([]string, 0, len(workPathSubstitutions)*2)
	for substitution, value := range workPathSubstitutions {
		oldNew = append(oldNew, substitution, value(w))
	}
	return strings.NewReplacer(oldNew...)
}

var pathSubstitutionRegexp = regexp.MustCompile(`{[^}]*}`)

var workPathSubstitutions = map[string]func(w *work.Work) string{
	"{title}": func(w *work.Work) string {
		return utils.FromPtr(w.Title, "Unknown")
	},
	"{id}": func(w *work.Work) string {
		return utils.FromPtrTransform(w.Id, utils.FormatUint64, "Unknown")
	},
	"{user}": func(w *work.Work) string {
		return utils.FromPtr(w.UserName, "Unknown")
	},
	"{user-id}": func(w *work.Work) string {
		return utils.FromPtrTransform(w.UserId, utils.FormatUint64, "Unknown")
	},
	"{type}": func(w *work.Work) string {
		switch utils.FromPtr(w.Kind, 255) {
		case work.KindIllust:
			return "Illustrations"
		case work.KindManga:
			return "Manga"
		case work.KindUgoira:
			return "Ugoira"
		case work.KindNovel:
			return "Novels"
		default:
			return "Unknown"
		}
	},
	"{restriction}": func(w *work.Work) string {
		switch utils.FromPtr(w.Restriction, 255) {
		case work.RestrictionNone:
			return "All Ages"
		case work.RestrictionR18:
			return "R-18"
		case work.RestrictionR18G:
			return "R-18G"
		default:
			return "Unknown"
		}
	},
	"{ai}": func(w *work.Work) string {
		if w.Ai == nil {
			return "Unknown"
		} else if *w.Ai {
			return "AI"
		} else {
			return "Human"
		}
	},
	"{original}": func(w *work.Work) string {
		if w.Original == nil {
			return "Unknown"
		} else if *w.Original {
			return "Original"
		} else {
			return "Not Original"
		}
	},
	"{series}": func(w *work.Work) string {
		return utils.FromPtr(w.SeriesTitle, "Unknown")
	},
	"{series-id}": func(w *work.Work) string {
		return utils.FromPtrTransform(w.SeriesId, utils.FormatUint64, "Unknown")
	},
}
