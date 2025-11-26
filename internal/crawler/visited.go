package crawler

import "sync"

type VisitedSet struct {
	urls map[string]struct{}
	mu   sync.Mutex
}

func NewVisitedSet() *VisitedSet {
	return &VisitedSet{
		urls: make(map[string]struct{}),
	}
}

func (v *VisitedSet) Add(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, ok := v.urls[url]; ok {
		return false
	}

	v.urls[url] = struct{}{}
	return true
}

func (v *VisitedSet) Len() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return len(v.urls)
}
