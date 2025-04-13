package queue

type ItemKind uint8

const (
	ItemKindArtwork ItemKind = iota
	ItemKindNovel
	ItemKindDefault = ItemKindArtwork

	ItemKindArtworkString = "artwork"
	ItemKindNovelString   = "novel"
	ItemKindDefaultString = ItemKindArtworkString
)

func ValidItemKindString(kind string) bool {
	return kind == ItemKindArtworkString || kind == ItemKindNovelString
}

func (kind ItemKind) String() string {
	switch kind {
	case ItemKindArtwork:
		return ItemKindArtworkString
	case ItemKindNovel:
		return ItemKindNovelString
	default:
		return ItemKindDefaultString
	}
}

func ItemKindFromString(kind string) ItemKind {
	switch kind {
	case ItemKindArtworkString:
		return ItemKindArtwork
	case ItemKindNovelString:
		return ItemKindNovel
	default:
		return ItemKindDefault
	}
}
