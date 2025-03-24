package work

type Restriction uint8

const (
	RestrictionNone Restriction = 0
	RestrictionR18  Restriction = 1
	RestrictionR18G Restriction = 2
)

func RestrictionOrDefault(restriction uint8) Restriction {
	if restriction <= 2 {
		return Restriction(restriction)
	}
	return RestrictionNone
}
