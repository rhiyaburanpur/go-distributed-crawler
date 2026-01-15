package main

import (
	"log"
	"sync"
	"time"

	"github.com/rhiyaburanpur/go-distributed-crawler/internal/client"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/crawler"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/storage"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/util"
)

func main() {
	const connStr = "postgres://crawler_user:crawler_pass@localhost:5432/crawler_metadata"
	const maxCrawls = 10
	const numWorkers = 5

	db, err := storage.NewPostgresDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	visited := crawler.NewVisitedSet(db)
	limiter := crawler.NewRateLimiter(1 * time.Second)
	var wg sync.WaitGroup

	urlCh := make(chan string, 100)
	resultCh := make(chan []string, numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, urlCh, resultCh, &wg, limiter, visited)
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
	log.Printf("Phase 3: Persistent Crawling completed. Total unique URLs in DB: %d\n", visited.Len())
}

func worker(id int, urlCh <-chan string, resultCh chan<- []string, wg *sync.WaitGroup, limiter *crawler.RateLimiter, visited *crawler.VisitedSet) {
	defer wg.Done()
	for currentURL := range urlCh {
		limiter.Wait(currentURL)

		var content string
		var err error
		delay := 1 * time.Second

		for i := 0; i < 3; i++ {
			content, err = client.Fetch(currentURL)
			if err == nil {
				break
			}
			time.Sleep(delay)
			delay *= 2
		}

		statusCode := 200
		errStr := ""
		if err != nil {
			statusCode = 0
			errStr = err.Error()
		}

		visited.UpdateStatus(currentURL, statusCode, errStr)

		// only extracting links if the fetch was successful
		if err == nil {
			newLinks := util.ExtractLinks(content, currentURL)
			resultCh <- newLinks
		} else {
			resultCh <- []string{}
		}
	}
}
