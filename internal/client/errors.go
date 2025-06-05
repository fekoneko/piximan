package client

type StatusError struct {
	error
	code int
}

func (e StatusError) Is(target error) bool {
	_, ok := target.(StatusError)
	return ok
}
