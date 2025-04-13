package queue

type Item struct {
	Id    uint64
	Kind  ItemKind
	Paths []string
}
