package work

import "github.com/fekoneko/piximan/internal/utils"

const (
	AiUnknownUint uint8 = 0
	AiNotAiUint   uint8 = 1
	AiIsAiUint    uint8 = 2
	AiDefaultUint       = AiUnknownUint
)

func AiFromUint(ai uint8) *bool {
	switch ai {
	case AiUnknownUint:
		return nil
	case AiNotAiUint:
		return utils.ToPtr(false)
	case AiIsAiUint:
		return utils.ToPtr(true)
	default:
		return AiFromUint(AiDefaultUint)
	}
}
