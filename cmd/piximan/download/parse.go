package download

import (
	"fmt"
	"strconv"
	"strings"
)

func parseIds(input string) ([]uint64, error) {
	parts := strings.Split(input, ",")
	ids := []uint64{}

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		id, err := strconv.ParseUint(trimmed, 10, 64)
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
	parts := strings.Split(input, ",")
	strs := []string{}

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			strs = append(strs, trimmed)
		}
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
