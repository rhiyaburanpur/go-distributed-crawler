# Go Distributed Crawler

A high-performance, scalable, and fault-tolerant web crawler built with Go. This project demonstrates distributed systems concepts including concurrent worker pools, rate limiting, persistent metadata storage, and cross-node coordination.

## Features

* Concurrent Worker Pool: Efficiently processes multiple URLs simultaneously using Go routines.

* Politeness and Rate Limiting: Per-host delay mechanism to prevent overwhelming target servers.

* Fault Tolerance: Implements exponential backoff for network retries and persists error metadata.

* Distributed Coordination: Uses Redis to manage a global work queue and atomic counters, allowing multiple instances to work together.

* Persistent Storage: Saves crawl results and metadata to PostgreSQL with automatic schema migration.

* Containerized Architecture: Fully orchestrated using Docker and Docker Compose for easy deployment and scaling.

![Project Demo](assets/demo.gif)

## Tech Stack

* Language: Go 1.25+

* Database: PostgreSQL (Metadata and URL persistence)

* Coordination: Redis (Distributed queue and state management)

* Infrastructure: Docker & Docker Compose

## Architecture

The system consists of three primary components:

1. Coordinator: Manages the global state, monitors the Redis work count, and feeds new discoveries back into the Redis queue.

2. Workers: Independent nodes that pull URLs from Redis, respect rate limits, fetch content, and report results back.

3. Storage Layer: A modular persistence layer handling SQL (Postgres) and NoSQL (Redis) operations.

## Project Structure

```
.
├── internal/
│   ├── client/     # HTTP fetch logic
│   ├── crawler/    # Politeness, visited set, and business logic
│   ├── storage/    # Postgres and Redis connection/logic
│   └── util/       # HTML parsing and link extraction
├── Dockerfile      # Multi-stage build for the Go application
├── docker-compose.yml # Service orchestration
├── main.go         # Application entry point and coordination loop
└── go.mod          # Dependency management
```

## Getting Started

### Prerequisites

* Docker Desktop installed and running.

* Go installed locally (for development/testing without containers).

### Installation and Running

1. Clone the repository:
```
git clone [https://github.com/rhiyaburanpur/go-distributed-crawler.git](https://github.com/rhiyaburanpur/go-distributed-crawler.git)
cd go-distributed-crawler
```

2. Launch the entire stack using Docker Compose:
```
docker-compose up --build
```

3. The crawler will start with a seed URL and begin discovering links. You can monitor the progress in the terminal logs.

### Scaling the Crawler

To demonstrate the distributed nature of the system, you can scale the number of worker nodes with a single command:
```
docker-compose up --build --scale app=3
```

This starts three separate containers running the Go application, all pulling from the same Redis queue and writing to the same Postgres database.

## Database Verification

To check the crawled data inside the Postgres container:
```
docker exec -it go-distributed-crawler-db-1 psql -U crawler_user -d crawler_metadata -c "SELECT url, status_code, crawled_at FROM crawled_urls LIMIT 20;"
```

## Development Workflow

This project was developed in six phases:

1. Basic HTTP Fetcher and Link Extractor.

2. Concurrent Worker Pool with Channels.

3. Persistent Storage with PostgreSQL.

4. Fault Tolerance and Rate Limiting.

5. Distributed Coordination with Redis.

6. Containerization and Orchestration.