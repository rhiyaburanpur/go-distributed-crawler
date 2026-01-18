package main

import (
	"context"
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
	const redisAddr = "localhost:6379"
	const maxCrawls = 15
	const numWorkers = 5

	// 1. Initialize Context and Redis
	ctx := context.Background()
	rdb := storage.NewRedisClient(redisAddr)
	defer rdb.Close()

	// 2. Clear previous Redis state for a fresh run
	rdb.Del(ctx, "crawler:work_count")

	// 3. Initialize Postgres
	db, err := storage.NewPostgresDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 4. Initialize Crawler Components
	visited := crawler.NewVisitedSet(db)
	limiter := crawler.NewRateLimiter(1 * time.Second)
	var wg sync.WaitGroup

	urlCh := make(chan string, 100)
	resultCh := make(chan []string, numWorkers)

	// 5. Start Workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, urlCh, resultCh, &wg, limiter, visited)
	}

	// 6. Seed the Crawl
	initialURL := "http://example.com"
	if visited.Add(initialURL) {
		rdb.Incr(ctx, "crawler:work_count")
		urlCh <- initialURL
	}

	// 7. Main Coordination Loop
	for {
		// Check global work count in Redis
		// .Int() handles conversion. If key doesn't exist, it returns 0 and an error.
		count, err := rdb.Get(ctx, "crawler:work_count").Int()

		// If count is 0 or key doesn't exist, we are done
		if err != nil || count <= 0 {
			break
		}

		// Wait for a worker to report results
		links := <-resultCh
		rdb.Decr(ctx, "crawler:work_count")

		for _, link := range links {
			if visited.Len() >= maxCrawls {
				break
			}
			if visited.Add(link) {
				rdb.Incr(ctx, "crawler:work_count")
				urlCh <- link
			}
		}
	}

	// 8. Graceful Shutdown
	close(urlCh)
	wg.Wait()
	log.Printf("Phase 5: Distributed Coordination started. Unique URLs in DB: %d\n", visited.Len())
}

func worker(id int, urlCh <-chan string, resultCh chan<- []string, wg *sync.WaitGroup, limiter *crawler.RateLimiter, visited *crawler.VisitedSet) {
	defer wg.Done()

	for currentURL := range urlCh {
		limiter.Wait(currentURL)

		var content string
		var err error
		delay := 1 * time.Second

		// Retry logic from Phase 4
		for i := 0; i < 3; i++ {
			content, err = client.Fetch(currentURL)
			if err == nil {
				break
			}
			time.Sleep(delay)
			delay *= 2
		}

		// Metadata reporting
		statusCode := 200
		errStr := ""
		if err != nil {
			statusCode = 0
			errStr = err.Error()
		}
		visited.UpdateStatus(currentURL, statusCode, errStr)

		// Report results back to coordinator
		if err == nil {
			newLinks := util.ExtractLinks(content, currentURL)
			resultCh <- newLinks
		} else {
			resultCh <- []string{}
		}
	}
}
