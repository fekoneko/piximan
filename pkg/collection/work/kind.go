package work

type Kind uint8

const (
	KindIllust Kind = 0
	KindManga  Kind = 1
)

const KindDefault = KindIllust

func KindFromUint(kind uint8) Kind {
	if kind <= 1 {
		return Kind(kind)
	}
	return KindDefault
}

func KindFromString(kind string) Kind {
	switch kind {
	case "illust":
		return KindIllust
	case "manga":
		return KindManga
	default:
		return KindDefault
	}
}

func (kind Kind) String() string {
	switch kind {
	case KindIllust:
		return "illust"
	case KindManga:
		return "manga"
	default:
		return KindDefault.String()
	}
}
