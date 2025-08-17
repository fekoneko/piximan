package fsext

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/internal/utils"
)

func InferIdsFromPattern(pattern string) (idPathMap *map[uint64][]string, errs []error) {
	pattern = filepath.Clean(pattern)
	r := inferIdRegexp(pattern)
	root, depth, err := inferIdWalkParams(pattern)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	idPathMap = utils.ToPtr(make(map[uint64][]string))
	inferIdWalk(root, r, 1, depth, idPathMap, &errs)

	return idPathMap, errs
}

func InferIdPatternValid(pattern string) error {
	patternIdIndex := strings.Index(pattern, "{id}")
	if patternIdIndex == -1 {
		return fmt.Errorf("pattern must contain {id}")
	}
	if strings.Contains(pattern[patternIdIndex+1:], "{id}") {
		return fmt.Errorf("pattern may not contain {id} twice")
	}
	return nil
}

// Returns weather the string could be an infer id pattern. This still doesn't mean it's valid.
func IsInferIdPattern(pattern string) bool {
	return inferIdSubstitutionRegexp.MatchString(pattern)
}

var inferIdSubstitutionRegexp = regexp.MustCompile(`{[^}]*}|\*`)

func inferIdRegexp(pattern string) *regexp.Regexp {
	prevIndex := 0
	builder := strings.Builder{}
	builder.WriteString(`^`)
	for _, bounds := range inferIdSubstitutionRegexp.FindAllStringIndex(pattern, -1) {
		plainText := pattern[prevIndex:bounds[0]]
		builder.WriteString(regexp.QuoteMeta(plainText))

		if pattern[bounds[0]:bounds[1]] == "{id}" {
			builder.WriteString(`([0-9]+)`)
		} else {
			builder.WriteString(`[^\\/]*`)
		}
		prevIndex = bounds[1]
	}
	plainText := pattern[prevIndex:]
	builder.WriteString(regexp.QuoteMeta(plainText))
	builder.WriteString(`$`)

	return regexp.MustCompile(builder.String())
}

func inferIdWalkParams(pattern string) (root string, depth int, err error) {
	bounds := inferIdSubstitutionRegexp.FindStringIndex(pattern)
	if bounds == nil {
		return "", 0, fmt.Errorf("pattern doesn't contain any substitutions")
	}
	separatorIndex := strings.LastIndex(pattern[:bounds[0]], string(os.PathSeparator))
	root = pattern[:max(separatorIndex, 0)]
	depth = strings.Count(pattern[separatorIndex+1:], string(os.PathSeparator)) + 1

	return root, depth, nil
}

func inferIdWalk(
	path string, r *regexp.Regexp, currentDepth int, depth int, idPathMap *map[uint64][]string, errs *[]error,
) {
	subEntries, err := os.ReadDir(filepath.Clean(path))
	if err != nil {
		*errs = append(*errs, err)
	}

	for _, subEntry := range subEntries {
		subName := subEntry.Name()
		subPath := filepath.Join(path, subName)

		if currentDepth < depth && subEntry.IsDir() {
			inferIdWalk(subPath, r, currentDepth+1, depth, idPathMap, errs)
		} else if currentDepth == depth {
			if matches := r.FindStringSubmatch(subPath); len(matches) > 0 {
				if id, err := strconv.ParseUint(matches[1], 10, 64); err == nil && id != 0 {
					(*idPathMap)[id] = append((*idPathMap)[id], subPath)
				} else {
					*errs = append(*errs, err)
				}
			}
		}
	}
}
