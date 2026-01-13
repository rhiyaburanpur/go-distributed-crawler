# Distributed Web Crawler: Phase 2 - Concurrency and Rate Limiting

## Project Overview

Phase 2 transitions the crawler from a sequential, single-threaded execution model to a concurrent, multi-worker architecture. This implementation leverages Go’s concurrency primitives -> Goroutines and Channels, to execute network I/O and HTML parsing in parallel. A critical addition in this phase is the Host-Based Rate Limiter, which ensures the system adheres to politeness protocols by throttling requests on a per-domain basis.

## Technical Specifications (Phase 2)

| Specification | Implementation | Purpose |
| :--- | :--- | :--- |
| **Execution Model** | **Goroutine Worker Pool** | **Parallelized network requests and content processing.** |
| **Communication** | **CSP (Communicating Sequential Processes)** | **Go channels facilitate data transfer between workers and the coordinator.** |
| **Politeness** | **Token Bucket/Timestamp Mapping** | **Enforces time delays between requests to the same host.** |
| **Synchronization** | **Sync Primitives** | **Mutexes and WaitGroups manage shared state and lifecycle orchestration.** |

## Architectural Design: Worker-Coordinator Pattern

The system has been refactored into two primary functional layers:

### 1. Central Coordinator

The coordinator, located in the main entry point, manages the global state of the crawl. It is responsible for:

* Maintaining the Visited Set to prevent redundant URL processing.

* Tracking in-flight operations via an atomic work counter.

* Orchestrating the lifecycle of the worker pool.

* Closing communication channels once the termination criteria (max crawls or exhausted queue) are met.

### 2. Concurrent Worker Pool

The worker pool consists of independent Goroutines that perform the following operations:

* Consuming URLs from the ingress channel.

* Querying the Rate Limiter to determine the necessary wait time for a specific host.

* Executing HTTP requests via the internal client.

* Parsing HTML and returning extracted links to the coordinator through the egress channel.

## Current Project Status

| Component | Status | Implementation Detail |
| :--- | :--- | :--- |
| **Coordinator Logic** | **Complete** | **Efficiently manages BFS state and synchronization.** |
| **Worker Pool** | **Complete** | **High-concurrency fetch/parse execution.** |
| **Rate Limiter** | **Complete** | **Thread-safe per-host access tracking using sync.Mutex.** |
| **Visited Set** | **Verified** | **Thread-safe set implementation for cycle detection.**| 
| **CI/CD Integration** | **Active** | **Automated validation of concurrent builds via GitHub Actions.** |

## Next Objectives

The subsequent phase (Phase 3) will focus on system persistence and fault tolerance. Key tasks include:

* Migrating the Visited Set from in-memory storage to a PostgreSQL database.

* Implementing a persistent work queue for crawl resume capability.

* Introducing exponential backoff for failed network requests.

## Deployment and Testing (Phase 2) 


To initialize the concurrent crawler and observe worker performance:

1. Clone the repository and navigate to the project root.

2. Initialize dependencies:
```
go mod download
```

3. Execute the crawler:
```
go run main.go
```

Note: The default configuration initializes 5 workers and a limit of 10 unique URLs. Progress and completion metrics will be reported via standard log output.
