package work

type Restriction uint8

const (
	RestrictionNone    Restriction = 0
	RestrictionR18     Restriction = 1
	RestrictionR18G    Restriction = 2
	RestrictionDefault             = RestrictionNone

	RestrictionNoneString    = "none"
	RestrictionR18String     = "R-18"
	RestrictionR18GString    = "R-18G"
	RestrictionDefaultString = RestrictionNoneString
)

func RestrictionFromUint(restriction uint8) Restriction {
	if restriction <= 2 {
		return Restriction(restriction)
	}
	return RestrictionDefault
}

func RestrictionFromString(restriction string) Restriction {
	switch restriction {
	case RestrictionNoneString:
		return RestrictionNone
	case RestrictionR18String:
		return RestrictionR18
	case RestrictionR18GString:
		return RestrictionR18G
	default:
		return RestrictionDefault
	}
}

func (restriction Restriction) String() string {
	switch restriction {
	case RestrictionNone:
		return RestrictionNoneString
	case RestrictionR18:
		return RestrictionR18String
	case RestrictionR18G:
		return RestrictionR18GString
	default:
		return RestrictionDefaultString
	}
}
