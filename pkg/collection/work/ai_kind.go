package work

type AiKind uint8

const (
	AiKindUnknown AiKind = 0
	AiKindNotAi   AiKind = 1
	AiKindIsAi    AiKind = 2
)

const AiKindDefault = AiKindUnknown

func AiKindFromUint(aiKind uint8) AiKind {
	if aiKind <= 2 {
		return AiKind(aiKind)
	}
	return AiKindDefault
}

func AiKindFromString(aiKind string) AiKind {
	switch aiKind {
	case "unknown":
		return AiKindUnknown
	case "not ai":
		return AiKindNotAi
	case "is ai":
		return AiKindIsAi
	default:
		return AiKindDefault
	}
}

func (aiKind AiKind) String() string {
	switch aiKind {
	case AiKindUnknown:
		return "unknown"
	case AiKindNotAi:
		return "not ai"
	case AiKindIsAi:
		return "is ai"
	default:
		return AiKindDefault.String()
	}
}
