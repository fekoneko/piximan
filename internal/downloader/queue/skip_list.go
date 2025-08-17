package queue

import (
	"github.com/fekoneko/piximan/internal/collection/work"
)

// Should be used to skip downloading works already present in the collection.
// SkipList differs from download rules in that it couples work ID with queue item kind
// and matches faster. Crawl tasks also may be aware of this list (e.g. skip fetching consequent
// bookmark pages if all found works on the current page were skipped).
type SkipList map[uint64]skipKind

func SkipListFromWorks(works []*work.Work) *SkipList {
	list := make(SkipList, len(works))

	for _, w := range works {
		if w.Id == nil || w.Kind == nil {
			continue
		} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga || *w.Kind == work.KindUgoira {
			list[*w.Id] = skipKindArtwork
		} else if *w.Kind == work.KindNovel {
			list[*w.Id] = skipKindNovel
		}
	}
	return &list
}

func SkipListFromMap[T any](m *map[uint64]T, kind ItemKind) *SkipList {
	list := make(SkipList, len(*m))
	for id := range *m {
		list[id] = skipKind(kind)
	}
	return &list
}

func (list *SkipList) Contains(id uint64, kind ItemKind) bool {
	found, ok := (*list)[id]
	return ok && found.match(kind)
}

type skipKind uint8

const (
	skipKindArtwork skipKind = iota
	skipKindNovel
	skipKindBoth
)

func (i skipKind) match(kind ItemKind) bool {
	return (i == skipKindArtwork && kind == ItemKindArtwork) ||
		(i == skipKindNovel && kind == ItemKindNovel) ||
		i == skipKindBoth
}
