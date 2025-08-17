package queue

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/logger"
)

// Should be used to skip downloading works already present in the collection.
// SkipList differs from download rules in that it couples work ID with queue item kind
// and matches faster. Crawl tasks also may be aware of this list (e.g. skip fetching consequent
// bookmark pages if all found works on the current page were skipped).
type SkipList struct {
	artworks map[uint64]bool
	novels   map[uint64]bool
}

func NewSkipList() *SkipList {
	return &SkipList{
		artworks: make(map[uint64]bool),
		novels:   make(map[uint64]bool),
	}
}

func SkipListFromWorks(works []*work.Work) *SkipList {
	list := NewSkipList()

	for _, w := range works {
		if w.Id != nil && w.Kind != nil {
			switch *w.Kind {
			case work.KindIllust, work.KindManga, work.KindUgoira:
				list.AddArtwork(*w.Id)
			case work.KindNovel:
				list.AddNovel(*w.Id)
			}
		} else {
			logger.Info("%v: %v", *w.Id, w.Kind)
		}
	}
	return list
}

func SkipListFromMap[T any](m *map[uint64]T, kind ItemKind) *SkipList {
	list := NewSkipList()
	for id := range *m {
		list.Add(id, kind)
	}
	return list
}

func (l *SkipList) Add(id uint64, kind ItemKind) {
	switch kind {
	case ItemKindArtwork:
		l.AddArtwork(id)
	case ItemKindNovel:
		l.AddNovel(id)
	}
}

func (l *SkipList) AddArtwork(id uint64) {
	l.artworks[id] = true
}

func (l *SkipList) AddNovel(id uint64) {
	l.novels[id] = true
}

func (l *SkipList) Remove(id uint64, kind ItemKind) {
	switch kind {
	case ItemKindArtwork:
		l.RemoveArtwork(id)
	case ItemKindNovel:
		l.RemoveNovel(id)
	}
}

func (l *SkipList) RemoveArtwork(id uint64) {
	delete(l.artworks, id)
}

func (l *SkipList) RemoveNovel(id uint64) {
	delete(l.novels, id)
}

func (l *SkipList) Contains(id uint64, kind ItemKind) bool {
	switch kind {
	case ItemKindArtwork:
		return l.ContainsArtwork(id)
	case ItemKindNovel:
		return l.ContainsNovel(id)
	default:
		return false
	}
}

func (l *SkipList) ContainsArtwork(id uint64) bool {
	return l.artworks[id]
}

func (l *SkipList) ContainsNovel(id uint64) bool {
	return l.novels[id]
}

func (l *SkipList) Len() int {
	return len(l.artworks) + len(l.novels)
}
