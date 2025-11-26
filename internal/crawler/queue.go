package crawler

import "sync"

type URLQueue struct {
	urls []string
	mu   sync.Mutex
}

func NewURLQueue() *URLQueue {
	return &URLQueue{
		urls: make([]string, 0),
	}
}

func (q *URLQueue) Enqueue(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.urls = append(q.urls, url)
}

func (q *URLQueue) Dequeue() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.urls) == 0 {
		return ""
	}

	url := q.urls[0]
	q.urls = q.urls[1:]
	return url
}

func (q *URLQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.urls) == 0
}
