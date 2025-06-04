package dto

type Response[T any] struct {
	Body T `json:"body"`
}

type ArtworkBookmarksBody struct {
	Works []BookmarkArtwork `json:"works"`
	Total uint64            `json:"total"`
}

type NovelBookmarksBody struct {
	Works []BookmarkNovel `json:"works"`
	Total uint64          `json:"total"`
}
