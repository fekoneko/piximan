package work

type AiKind uint8

const (
	AiKindUnknown AiKind = 0
	AiKindNotAi   AiKind = 1
	AiKindIsAi    AiKind = 2
)

const AiKindDefault = AiKindUnknown

func AiKindOrDefault(aiKind uint8) AiKind {
	if aiKind <= 2 {
		return AiKind(aiKind)
	}
	return AiKindDefault
}
