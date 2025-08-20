package download

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: move to internal/promptuiext and add more validation functions?

func parseIds(input string) ([]uint64, error) {
	strs := parseStrings(input)
	ids := []uint64{}

	for _, str := range strs {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("IDs must be a comma-separated list of numbers")
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no IDs provided")
	}
	return ids, nil
}

func parseStrings(input string) []string {
	strs := []string{}
	strStart := 0

	for i, r := range input {
		if (r == ',' || r == ';' || r == '、' || r == '；') && (i == 0 || input[i-1] != '\\') {
			if strStart < i {
				strs = append(strs, input[strStart:i])
			}
			strStart = i + 1
		}
	}
	if strStart < len(input) {
		strs = append(strs, input[strStart:])
	}

	return strs
}

func parseRange(input string) (from *uint64, to *uint64, err error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, nil, nil
	}

	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("the range must contain exactly one ':'")
	}

	trimmedFromString := strings.TrimSpace(parts[0])
	if trimmedFromString != "" {
		parsed, err := strconv.ParseUint(trimmedFromString, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'from' value: %v", from)
		}
		from = &parsed
	}

	trimmedToString := strings.TrimSpace(parts[1])
	if trimmedToString != "" {
		parsed, err := strconv.ParseUint(trimmedToString, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'to' value: %v", to)
		}
		to = &parsed
	}

	if from != nil && to != nil && *from >= *to {
		return nil, nil, fmt.Errorf("'from' value must be less than 'to' value")
	}

	return from, to, nil
}
