# Go Web Crawler Core (Phase 1)

## Project Goal

The primary goal of this phase is to build the sequential core for an efficient web crawler. This involves implementing the foundational data structures and logic in Go to perform a Breadth-First Search (BFS) traversal of links, ensure URL uniqueness, and handle basic HTTP fetching and HTML parsing.

## Key Technologies (Phase 1)

| Category | Purpose | Technology/Tool | Rationale/Detail |
| :--- | :--- | :--- | :--- |
| **Technology** | **Language** | Go (Golang) | Used for its high **performance** and strong focus on data structure implementation. |
| **Technology** | **Parsing** | `golang.org/x/net/html` | Chosen for robust **HTML parsing** and reliable link extraction. |
| **DevOps** | **CI/CD** | GitHub Actions | Handles **Continuous Integration (CI)** to automatically validate build and test success. |

## Phase 1 Architecture: Monolithic Core

In this initial phase, the system operates as a single, sequential application (main.go). All logic—fetching, parsing, and link management—is executed in a single loop, acting as a core proof-of-concept for the future distributed system.

**Core Components Implemented in Phase 1:**

* URL Queue: Implements the BFS order of operations.

* Visited Set: A thread-safe data structure to enforce URL uniqueness and prevent cycles.

* HTTP Client: Custom configuration for fetching web page content.

* HTML Parser: Logic to extract absolute URLs from fetched HTML content.

* Project Status: Phase 1 Complete

This phase successfully establishes the core logic and data structures required for web crawling:

| Component | Status | Location | Focus |
| :--- | :--- | :--- | :--- |
| **Core BFS Logic** | Complete | `main.go` | Sequential Breadth-First Search loop. |
| **Data Structures** | Complete | `internal/crawler/` | Thread-safe URLQueue and VisitedSet. |
| **HTTP Client** | Complete | `internal/client/http.go` | Custom client with sensible timeouts. |
| **HTML Parser** | Complete | `internal/util/parser.go` | Link extraction and absolute URL resolution. |
| **CI Workflow** | Implemented | `.github/workflows/go.yml` | Automatic build and test validation on commit. |

** Next Major Steps :**
 Phase 2 will introduce Goroutines, Channels, and the Politeness RateLimiter to transition this core into a multi-threaded, concurrent crawler.

 Getting Started (Phase 1 Test)

To run the single-threaded crawler locally and confirm the core logic:

Clone the repository:

```
git clone [your_repository_url]
cd go-distributed-crawler
```

Run the Phase 1 entry point:
```
go run main.go
```


(Note: This executes the sequential loop, fetches the seed URL, and stops after 10 unique links.)