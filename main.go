package main

import (
	"log"
	"sync"
	"time"

	"github.com/rhiyaburanpur/go-distributed-crawler/internal/client"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/crawler"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/util"
)

func main() {

	const maxCrawls = 10
	const numWorkers = 5

	visited := crawler.NewVisitedSet()
	limiter := crawler.NewRateLimiter(1 * time.Second)
	var wg sync.WaitGroup

	urlCh := make(chan string, 100)
	resultCh := make(chan []string, numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, urlCh, resultCh, &wg, limiter)
	}

	workCount := 0
	initialURL := "http://example.com"

	if visited.Add(initialURL) {
		workCount++
		urlCh <- initialURL
	}

	for workCount > 0 {
		links := <-resultCh
		workCount--

		for _, link := range links {
			if visited.Len() >= maxCrawls {
				break
			}
			if visited.Add(link) {
				workCount++
				urlCh <- link
			}
		}
	}
	close(urlCh)
	wg.Wait()
	log.Printf("Phase 2: Concurrent Crawling completed. Visited %d unique URLs.\n", visited.Len())
}

func worker(id int, urlCh <-chan string, resultCh chan<- []string, wg *sync.WaitGroup, limiter *crawler.RateLimiter) {
	defer wg.Done()

	for currentURL := range urlCh {
		limiter.Wait(currentURL)
		content, err := client.Fetch(currentURL)
		if err != nil {
			continue
		}
		newLinks := util.ExtractLinks(content, currentURL)
		resultCh <- newLinks
	}
}
