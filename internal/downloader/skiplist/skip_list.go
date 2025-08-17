package skiplist

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
)

// Should be used to skip downloading works already present in the collection.
// SkipList differs from download rules in that it couples work ID with queue item kind
// and matches faster. Crawl tasks also may be aware of this list (e.g. skip fetching consequent
// bookmark pages if all found works on the current page were skipped).
type SkipList struct {
	artworks map[uint64]bool
	novels   map[uint64]bool
}

func New() *SkipList {
	return &SkipList{
		artworks: make(map[uint64]bool),
		novels:   make(map[uint64]bool),
	}
}

func (l *SkipList) AddWork(w *work.Work) {
	if w.Id != nil && w.Kind != nil {
		switch *w.Kind {
		case work.KindIllust, work.KindManga, work.KindUgoira:
			l.AddArtwork(*w.Id)
		case work.KindNovel:
			l.AddNovel(*w.Id)
		}
	}
}

func (l *SkipList) Add(id uint64, kind queue.ItemKind) {
	switch kind {
	case queue.ItemKindArtwork:
		l.AddArtwork(id)
	case queue.ItemKindNovel:
		l.AddNovel(id)
	}
}

func (l *SkipList) AddArtwork(id uint64) {
	l.artworks[id] = true
}

func (l *SkipList) AddNovel(id uint64) {
	l.novels[id] = true
}

func (l *SkipList) Remove(id uint64, kind queue.ItemKind) {
	switch kind {
	case queue.ItemKindArtwork:
		l.RemoveArtwork(id)
	case queue.ItemKindNovel:
		l.RemoveNovel(id)
	}
}

func (l *SkipList) RemoveArtwork(id uint64) {
	delete(l.artworks, id)
}

func (l *SkipList) RemoveNovel(id uint64) {
	delete(l.novels, id)
}

func (l *SkipList) Contains(id uint64, kind queue.ItemKind) bool {
	switch kind {
	case queue.ItemKindArtwork:
		return l.ContainsArtwork(id)
	case queue.ItemKindNovel:
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
