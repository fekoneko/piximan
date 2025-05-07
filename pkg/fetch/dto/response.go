package dto

type Response[T any] struct {
	Body T `json:"body"`
}

type BookmarkArtworksBody struct {
	Works []BookmarkArtwork `json:"works"`
}
