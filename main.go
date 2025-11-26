package main

import (
	"fmt"
	"log"
	"time"

	// IMPORTANT: Use your specific module path
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/client"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/crawler"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/util"
)

func main() {
	log.Println("--- Phase 1: Single-Threaded Crawler Starting ---")

	// 1. Initialize the required DSA structures for BFS and uniqueness
	queue := crawler.NewURLQueue()
	visited := crawler.NewVisitedSet()

	// Define the initial seed URL
	seedURL := "http://example.com"

	// Enqueue the first URL and mark it as visited
	queue.Enqueue(seedURL)
	visited.Add(seedURL)

	count := 0
	// Limit the crawl for testing
	maxCrawls := 10

	// 2. Start the single-threaded crawl loop (BFS)
	for !queue.IsEmpty() && count < maxCrawls {
		// Dequeue the next URL to visit
		currentURL := queue.Dequeue()
		count++

		log.Printf("Crawling (%d/%d): %s\n", count, maxCrawls, currentURL)

		// 3. Fetch the content using the custom HTTP Client
		content, err := client.Fetch(currentURL)
		if err != nil {
			log.Printf("ERROR fetching %s: %v", currentURL, err)
			continue
		}

		// Simulate politeness delay (good practice, even in sequential mode)
		time.Sleep(1 * time.Second)

		// 4. Extract and process new links using the Parser
		newLinks := util.ExtractLinks(content, currentURL)

		for _, link := range newLinks {
			// 5. Enforce uniqueness using the Visited Set (Hash Set logic)
			if visited.Add(link) {
				// If the link is new, add it to the queue for next-level BFS exploration
				queue.Enqueue(link)
				fmt.Printf("  -> Found new link: %s\n", link)
			}
		}
	}

	log.Println("\n--- Phase 1: Single-Threaded Crawl Finished ---")
	log.Printf("Total unique URLs processed: %d\n", count)
}
