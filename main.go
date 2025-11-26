package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rhiyaburanpur/go-distributed-crawler/internal/client"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/crawler"
	"github.com/rhiyaburanpur/go-distributed-crawler/internal/util"
)

func main() {
	log.Println("Phase 1: Single-Threaded Crawler Starting")

	queue := crawler.NewURLQueue()
	visited := crawler.NewVisitedSet()

	seedURL := "http://example.com"
	queue.Enqueue(seedURL)
	visited.Add(seedURL)

	count := 0
	maxCrawls := 10

	for !queue.IsEmpty() && count < maxCrawls {
		currentURL := queue.Dequeue()
		count++

		log.Printf("Crawling (%d/%d): %s\n", count, maxCrawls, currentURL)

		content, err := client.Fetch(currentURL)
		if err != nil {
			log.Printf("ERROR fetching %s: %v", currentURL, err)
			continue
		}

		time.Sleep(1 * time.Second)

		newLinks := util.ExtractLinks(content, currentURL)

		for _, link := range newLinks {
			if visited.Add(link) {
				queue.Enqueue(link)
				fmt.Printf("  -> Found new link: %s\n", link)
			}
		}
	}

	log.Println("\nPhase 1: Single-Threaded Crawl Finished")
	log.Printf("Total unique URLs processed: %d\n", count)
}
