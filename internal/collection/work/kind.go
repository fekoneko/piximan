package work

type Kind uint8

const (
	KindIllust  Kind = 0
	KindManga   Kind = 1
	KindUgoira  Kind = 2
	KindNovel   Kind = 3
	KindDefault      = KindIllust

	KindIllustString  = "illust"
	KindMangaString   = "manga"
	KindUgoiraString  = "ugoira"
	KindNovelString   = "novel"
	KindDefaultString = KindIllustString
)

func KindFromUint(kind uint8) Kind {
	if kind <= 3 {
		return Kind(kind)
	}
	return KindDefault
}

func KindFromString(kind string) Kind {
	switch kind {
	case KindIllustString:
		return KindIllust
	case KindMangaString:
		return KindManga
	case KindUgoiraString:
		return KindUgoira
	case KindNovelString:
		return KindNovel
	default:
		return KindDefault
	}
}

func (kind Kind) String() string {
	switch kind {
	case KindIllust:
		return KindIllustString
	case KindManga:
		return KindMangaString
	case KindUgoira:
		return KindUgoiraString
	case KindNovel:
		return KindNovelString
	default:
		return KindDefault.String()
	}
}
