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

func AiKindFromBool(aiKind *bool) AiKind {
	if aiKind == nil {
		return AiKindUnknown
	}
	if *aiKind {
		return AiKindIsAi
	}
	return AiKindNotAi
}

func (aiKind AiKind) Bool() *bool {
	switch aiKind {
	case AiKindUnknown:
		return nil
	case AiKindNotAi:
		return new(bool)
	case AiKindIsAi:
		result := new(bool)
		*result = true
		return result
	default:
		return AiKindDefault.Bool()
	}
}
