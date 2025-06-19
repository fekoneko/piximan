package fsext

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var inferPatternReplacer = strings.NewReplacer(
	"\\", "\\\\",
	"[", "\\[",
	"]", "\\]",
	"?", "\\?",
)

// TODO: refactor this abomination
func InferIdsFromWorkPath(pattern string) (idPathMap *map[uint64][]string, err error) {
	pattern, idIndex, err := FormatInferIdPath(pattern)
	if err != nil {
		return nil, err
	}
	matches, err := filepath.Glob(strings.ReplaceAll(pattern, "{id}", "*"))
	if err != nil {
		return nil, err
	}

	separator := string(os.PathSeparator)
	slashesAfterId := strings.Count(pattern[idIndex:], separator)
	end := strings.Index(pattern[idIndex:], separator)
	if end != -1 {
		end += idIndex
	} else {
		end = len(pattern)
	}
	start := strings.LastIndex(pattern[:idIndex], separator) + 1
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

func InferIdPathValid(pattern string) error {
	_, _, err := FormatInferIdPath(pattern)
	return err
}

// Replace all {substitutions} with *, but keep {id}
func FormatInferIdPath(pattern string) (formatted string, idIndex int, err error) {
	pattern = inferPatternReplacer.Replace(pattern)
	pattern = patternRegex.ReplaceAllStringFunc(pattern, func(s string) string {
		if s == "{id}" {
			return s
		}
		return "*"
	})

	patternIdIndex := strings.Index(pattern, "{id}")
	if patternIdIndex == -1 {
		return "", 0, fmt.Errorf("pattern must contain {id}")
	}
	if strings.Contains(pattern[patternIdIndex+1:], "{id}") {
		return "", 0, fmt.Errorf("pattern may not contain {id} twice")
	}
	if (patternIdIndex >= 1 && pattern[patternIdIndex-1] == '*') ||
		(patternIdIndex < len(pattern)-len("{id}") && pattern[patternIdIndex+len("{id}")] == '*') {
		return "", 0, fmt.Errorf("pattern may not contain another substitution directly before or directly after {id}")
	}
	return pattern, patternIdIndex, nil
}
