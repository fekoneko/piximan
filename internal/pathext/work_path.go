package pathext

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/utils"
)

func FormatWorkPath(pattern string, w *work.Work) (string, error) {
	replacer := strings.NewReplacer(
		"{title}", w.Title,
		"{id}", strconv.FormatUint(w.Id, 10),
		"{user}", w.UserName,
		"{userid}", strconv.FormatUint(w.UserId, 10),
		"{restrict}", utils.If(
			w.Restriction == work.RestrictionNone, "all-ages", w.Restriction.String(),
		),
	)

	path, err := filepath.Abs(pattern)
	if err != nil {
		return "", err
	}
	sections := strings.Split(path, string(filepath.Separator))
	if len(sections) != 0 && len(sections[0]) == 0 {
		sections[0] = string(filepath.Separator)
	}

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

var inferPatternReplacer = strings.NewReplacer(
	"\\", "\\\\",
	"[", "\\[",
	"]", "\\]",
	"?", "\\?",
)

// TODO: refactor this abomination
func InferIdsFromWorkPath(pattern string) (*map[uint64][]string, error) {
	pattern = inferPatternReplacer.Replace(pattern)
	patternIdIndex := strings.Index(pattern, "{id}")
	if patternIdIndex == -1 {
		return nil, fmt.Errorf("pattern must contain {id}")
	}
	if strings.Contains(pattern[patternIdIndex+1:], "{id}") {
		return nil, fmt.Errorf("pattern may not contain {id} twice")
	}
	if (patternIdIndex >= 1 && pattern[patternIdIndex-1] == '*') ||
		(patternIdIndex < len(pattern)-len("{id}") && pattern[patternIdIndex+len("{id}")] == '*') {
		return nil, fmt.Errorf("pattern may not contain * directly before or directly after {id}")
	}

	matches, err := filepath.Glob(strings.ReplaceAll(pattern, "{id}", "*"))
	if err != nil {
		return nil, err
	}

	separator := string(os.PathSeparator)
	slashesAfterId := strings.Count(pattern[patternIdIndex:], separator)
	end := strings.Index(pattern[patternIdIndex:], separator)
	if end != -1 {
		end += patternIdIndex
	} else {
		end = len(pattern)
	}
	start := strings.LastIndex(pattern[:patternIdIndex], separator) + 1
	patternIdSection := pattern[start:end]
	result := make(map[uint64][]string)

	for _, match := range matches {
		matchIdSection := match[:]
		for i := 0; i < slashesAfterId; i++ {
			slashIndex := strings.LastIndex(matchIdSection, string(os.PathSeparator))
			if slashIndex == -1 {
				break
			}
			matchIdSection = matchIdSection[:slashIndex]
		}
		matchIdSection = filepath.Base(matchIdSection)

		m := []rune(matchIdSection)
		p := []rune(patternIdSection)
		idRunes := []rune{}

		for mi, pi := 0, 0; mi < len(m) && pi < len(p)-len("{id}")+1; mi, pi = mi+1, pi+1 {
			if p[pi] == '\\' {
				pi++
			} else if p[pi] == '*' {
				for ; mi < len(m) && (p[pi+1] == '\\' && m[mi] != p[pi+2] ||
					p[pi+1] != '\\' && m[mi] != p[pi+1]); mi++ {
				}
				mi--
			} else if p[pi] == '{' && p[pi+1] == 'i' && p[pi+2] == 'd' && p[pi+3] == '}' {
				pi += len("{id}")
				for ; mi < len(m) && (pi >= len(p) || (p[pi] == '\\' && pi+1 < len(p) &&
					m[mi] != p[pi+1] || p[pi] != '\\' && m[mi] != p[pi])); mi++ {
					idRunes = append(idRunes, m[mi])
				}
				break
			}
		}

		if id, err := strconv.ParseUint(string(idRunes), 10, 64); err == nil {
			result[id] = append(result[id], match)
		}
	}

	return &result, nil
}
