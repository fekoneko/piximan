package fsext

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fekoneko/piximan/internal/utils"
	"github.com/fekoneko/piximan/internal/work"
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

	replacer := getWorkPathReplacer(w)
	for i, section := range sections {
		if i == 0 || section == "." || section == ".." {
			continue
		}
		filename := replacer.Replace(section)
		sections[i] = ToValidFilename(filename)
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

var patternRegex = regexp.MustCompile(`{[^}]*}`)

func WorkPathValid(pattern string) error {
	for _, match := range patternRegex.FindAllString(pattern, -1) {
		if _, ok := workPathSubstitutions[match]; !ok {
			return fmt.Errorf("pattern contains unknown substitution %q", match)
		}
	}
	return nil
}

func getWorkPathReplacer(w *work.Work) *strings.Replacer {
	oldNew := make([]string, 0, len(workPathSubstitutions)*2)
	for substitution, value := range workPathSubstitutions {
		oldNew = append(oldNew, substitution, value(w))
	}

	return strings.NewReplacer(oldNew...)
}

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
		switch utils.FromPtr(w.AiKind, 255) {
		case work.AiKindNotAi:
			return "Human"
		case work.AiKindIsAi:
			return "AI"
		default:
			return "Unknown"
		}
	},
	"{original}": func(w *work.Work) string {
		if w.Original == nil {
			return "Unknown"
		}
		if *w.Original {
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
