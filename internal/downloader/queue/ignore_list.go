package queue

import (
	"github.com/fekoneko/piximan/internal/collection/work"
)

// Should be used to ignore downloading works already present in the collection.
// IgnoreList differs from download rules in that it couples work ID with queue item kind
// and matches faster. Crawl tasks also may be aware of this list (e.g. skip fetching consequent
// bookmark pages if all found works on the current page were ignored).
type IgnoreList map[uint64]ignoreKind

func IgnoreListFromWorks(works []*work.Work) *IgnoreList {
	list := make(IgnoreList, len(works))

	for _, w := range works {
		if w.Id == nil || w.Kind == nil {
			continue
		} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga || *w.Kind == work.KindUgoira {
			list[*w.Id] = ignoreKindArtwork
		} else if *w.Kind == work.KindNovel {
			list[*w.Id] = ignoreKindNovel
		}
	}
	return &list
}

func (list *IgnoreList) Contains(id uint64, kind ItemKind) bool {
	found, ok := (*list)[id]
	return ok && found.match(kind)
}

type ignoreKind uint8

const (
	ignoreKindArtwork ignoreKind = iota // TODO: remove UPPER_CASE everywhere
	ignoreKindNovel
	ignoreKindBoth
)

func (i ignoreKind) match(kind ItemKind) bool {
	return (i == ignoreKindArtwork && kind == ItemKindArtwork) ||
		(i == ignoreKindNovel && kind == ItemKindNovel) ||
		i == ignoreKindBoth
}
