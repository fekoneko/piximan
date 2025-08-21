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
	time, _ := ParseLocalTimePtrStrict(s)
	return time
}

// Parse time from RFC3339 string and return it in local time zone. Return error if parsing fails.
func ParseLocalTimePtrStrict(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	time, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil, err
	}
	localTime := time.Local()
	return &localTime, nil
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

func NoneDefined(values ...any) bool {
	for _, value := range values {
		if reflect.ValueOf(value).Pointer() != 0 {
			return false
		}
	}
	return true
}

func MapFindValue[K comparable, V comparable](m map[K]V, value V) (key K, ok bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return key, false
}

// Copies pointer value and returns a pointer to it.
func Copy[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	value := *ptr
	return &value
}

// Returns singular if n is 1. Use to pluralize words.
func IfPlural(n int, singular string, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

// Returns "s" if n is not 1. Use to pluralize words.
func Plural(n int) string {
	return IfPlural(n, "", "s")
}
