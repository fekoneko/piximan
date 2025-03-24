package work

type Kind uint8

const (
	KindIllust Kind = 0
	KindManga  Kind = 1
)

func KindOrDefault(kind uint8) Kind {
	if kind <= 1 {
		return Kind(kind)
	}
	return KindIllust
}
