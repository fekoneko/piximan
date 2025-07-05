package work

import "github.com/fekoneko/piximan/internal/utils"

type AiKind uint8

const (
	AiKindUnknown AiKind = 0
	AiKindNotAi   AiKind = 1
	AiKindIsAi    AiKind = 2
	AiKindDefault        = AiKindUnknown
)

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
		return utils.ToPtr(false)
	case AiKindIsAi:
		return utils.ToPtr(true)
	default:
		return AiKindDefault.Bool()
	}
}
