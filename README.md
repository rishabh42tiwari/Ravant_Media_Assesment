# Ravant Media Assessment

This repository contains my solutions for the Ravant Media technical assessment. The focus is on designing scalable, fault-tolerant systems using Go (Golang) — especially around concurrency, real-time ingestion, filtering at scale, and resilient API design.

---

## 1. Golang Concurrency in Action

This module demonstrates how to efficiently handle a large number of file operations using Go's concurrency features — goroutines, channels, and a worker pool.

### Real-World Scenario

In one of our production batch jobs, we had to process around 100 files fetched from a feed hub. For each file:

- Read and validate contents
- Insert into MongoDB

I built a worker pool using goroutines:
- A buffered channel pushed file paths into the pool
- Workers picked up tasks and ran the read-validate-write logic

The job was scheduled monthly or quarterly via OpenShift CronJobs. This setup offered:

- Efficient CPU and memory usage
- Safe parallel processing
- Predictable and scalable performance

---

## 2. Real-Time Ingestion Pipeline

The ingestion system is designed to handle 50,000+ entries per minute — balancing throughput, latency, and resilience.

### Design Highlights

- **Concurrency**: Go-based ingestion APIs use goroutines and channels to handle requests in parallel
- **Decoupling**: Data is sent to a message queue like Kafka or NATS for asynchronous processing
- **Streaming**: Worker pools and batch inserts reduce write latency
- **Fault Tolerance**:
  - Retries with exponential backoff
  - Dead letter queues for unrecoverable events
  - Durable queues to allow safe replay
- **Monitoring**: Prometheus + Grafana for system metrics
- **Caching**: Redis caches hot data and reduces database hits

---

## 3. Time-Based Queries with MongoDB

For storing time-series-like data, MongoDB works well due to its flexible schema and powerful query capabilities.

### Schema Example

```json
{
  "deviceId": "abc-123",
  "A": 15,
  "B": 3,
  "createdAt": "2025-07-10T12:00:00Z"
}
```

# Scalable Filtering & Resilient API System

This project demonstrates how to build a **scalable, resilient system** that supports:

- User-defined filtering on large datasets
- High-performance REST APIs
- Fault-tolerant architecture for real-world reliability

---

## 4. User-Defined Filters at Scale

### Use Case
Users can define dynamic filters like:  
`"Show items where A > 10 and B < 5 for the last 20 records"`

### Solution Strategy

- **Flexible Schema**  
  Data is stored in MongoDB with fields like `A`, `B`, and `createdAt`.

- **Safe Filter Parsing**  
  User-defined filters (e.g., `"A > 10 AND B < 5"`) are parsed using a secure, custom-built parser and translated into MongoDB queries.

- **Efficient Querying**
  - Use `.sort({ createdAt: -1 }).limit(N)` to fetch the latest N records.
  - Apply parsed filters in the query itself for performance.

- **Concurrent Execution**
  - Use Goroutines and a worker pool to apply filters concurrently across thousands of datasets.
  - Control concurrency level to prevent database overload.

### Performance Tips

- Index `createdAt` and frequently queried fields (`A`, `B`, etc.)
- Cache repeated filter results to avoid redundant computation
- Use pagination or time-based slicing for handling large data volumes

---

## 5. API Performance & Scalability

To handle thousands of users with low response times, the system includes:

### Key Strategies

- **Caching**
  - Use Redis to cache frequent queries (charts, filters) and reduce load on MongoDB

- **Pagination**
  - All endpoints return paginated responses to avoid large payloads

- **Rate Limiting**
  - Token bucket or fixed window algorithms enforce per-user/IP limits (e.g., using Redis or Go's `rate` package)

- **Concurrent Handling**
  - Go's Goroutines and non-blocking I/O efficiently handle concurrent requests

- **Async Processing**
  - Heavy operations like analytics or notifications are handled via background workers or message queues

---

## 6. System Resilience & Fault Tolerance

The architecture is built to remain stable under service failures or network issues.

### Fault-Tolerant Patterns

- **Retry with Backoff**
  - Failed operations are retried with exponential backoff to reduce pressure on dependent services

- **Timeouts & Context Cancellation**
  - External calls are protected with timeouts using `context.WithTimeout`

- **Circuit Breakers**
  - Use libraries like `sony/gobreaker` to stop calling failing services temporarily and auto-recover

- **Fallbacks**
  - Serve cached or default responses during downtime

- **Dead Letter Queues (DLQ)**
  - Store failed or unprocessed data for later retries

- **Graceful Degradation**
  - Critical features stay online even if non-essential services fail

---

## Tech Stack

- **Backend**: Golang, Fiber
- **Database**: MongoDB
- **Cache**: Redis
- **Concurrency**: Goroutines, Worker Pools
- **Rate Limiting**: Redis or Go `rate` package
- **Monitoring (Optional)**: Prometheus, Grafana
- **Queue (Optional)**: NATS, RabbitMQ, or Kafka

---