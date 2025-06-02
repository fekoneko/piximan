package download

import (
	"fmt"
	"strconv"
	"strings"
)

func parseIds(idsString string) ([]uint64, error) {
	idSubstrs := strings.Split(idsString, ",")
	ids := []uint64{}

	for _, idSubstr := range idSubstrs {
		trimmed := strings.TrimSpace(idSubstr)
		if trimmed == "" {
			continue
		}
		id, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no IDs provided")
	}

	return ids, nil
}

func parseRange(rangeString string) (*uint64, *uint64, error) {
	trimmed := strings.TrimSpace(rangeString)
	if trimmed == "" {
		return nil, nil, nil
	}

	parts := strings.Split(rangeString, ":")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("the range must contain exactly one ':'")
	}

	var fromOffset *uint64
	var toOffset *uint64

	trimmedFromString := strings.TrimSpace(parts[0])
	if trimmedFromString != "" {
		parsed, err := strconv.ParseUint(trimmedFromString, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'from' value: %v", fromOffset)
		}
		fromOffset = &parsed
	}

	trimmedToString := strings.TrimSpace(parts[1])
	if trimmedToString != "" {
		parsed, err := strconv.ParseUint(trimmedToString, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'to' value: %v", toOffset)
		}
		toOffset = &parsed
	}

	if fromOffset != nil && toOffset != nil && *fromOffset >= *toOffset {
		return nil, nil, fmt.Errorf("'from' value must be less than 'to' value")
	}

	return fromOffset, toOffset, nil
}
