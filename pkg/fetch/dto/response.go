package dto

type Response[T any] struct {
	Body T `json:"body"`
}

type ArtworkBookmarksBody struct {
	Works []BookmarkArtwork `json:"works"`
}

type NovelBookmarksBody struct {
	Works []BookmarkNovel `json:"works"`
}
