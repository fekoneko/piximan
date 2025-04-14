package queue

type Item struct {
	Id       uint64
	Kind     ItemKind
	OnlyMeta bool
	Paths    []string
}
