package utils

import (
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

func FormatUTCTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.UTC().Format(time.RFC3339)
	return &formatted
}
