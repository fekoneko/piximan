package util

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
