package work

type Restriction uint8

const (
	RestrictionNone Restriction = 0
	RestrictionR18  Restriction = 1
	RestrictionR18G Restriction = 2
)

const RestrictionDefault = RestrictionNone

func RestrictionFromUint(restriction uint8) Restriction {
	if restriction <= 2 {
		return Restriction(restriction)
	}
	return RestrictionDefault
}

func RestrictionFromString(restriction string) Restriction {
	switch restriction {
	case "none":
		return RestrictionNone
	case "r18":
		return RestrictionR18
	case "r18g":
		return RestrictionR18G
	default:
		return RestrictionDefault
	}
}

func (restriction Restriction) String() string {
	switch restriction {
	case RestrictionNone:
		return "none"
	case RestrictionR18:
		return "r18"
	case RestrictionR18G:
		return "r18g"
	default:
		return RestrictionDefault.String()
	}
}
