package downloader

type ApiResponse[T any] struct {
	Body T `json:"body"`
}
