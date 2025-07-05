package utils

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func If[T any](condition bool, value T, defaultValue T) T {
	if condition {
		return value
	}
	return defaultValue
}

func FromPtr[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func FromPtrTransform[T any, U any](ptr *T, transform func(T) U, defaultValue U) U {
	if ptr == nil {
		return defaultValue
	}
	return transform(*ptr)
}

func MapPtr[T any, U any](ptr *T, transform func(T) U) *U {
	if ptr == nil {
		return nil
	}
	return ToPtr(transform(*ptr))
}

func ToPtr[T any](value T) *T {
	return &value
}

func ValidateNumber(message string) func(string) error {
	return func(s string) error {
		if _, err := strconv.ParseUint(s, 10, 64); err != nil {
			return errors.New(message)
		}
		return nil
	}
}

func FormatUint64(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func ParseUint64Ptr(s *string) *uint64 {
	if s == nil {
		return nil
	}
	value, err := strconv.ParseUint(*s, 10, 64)
	if err != nil {
		return nil
	}
	return &value
}

// Parse time from RFC3339 string and return it in local time zone.
func ParseLocalTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	time, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	localTime := time.Local()
	return &localTime
}

// Convert time to UTC and format it as RFC3339.
func FormatUTCTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.UTC().Format(time.RFC3339)
	return &formatted
}

func ExactlyOneDefined(values ...any) bool {
	defined := false
	for _, value := range values {
		if reflect.ValueOf(value).Pointer() != 0 {
			if defined {
				return false
			}
			defined = true
		}
	}
	return defined
}
