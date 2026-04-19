# High-Performance Rate Limiter in Go

A blazingly fast, concurrent-safe rate limiting service written in Go, implementing the **Token Bucket** algorithm. Currently optimized for single-node, high-throughput environments with a roadmap for distributed scaling.

## Features

- **Token Bucket Algorithm:** Smooth and predictable rate limiting.
- **Concurrent-Safe:** Built from the ground up to handle massive concurrent requests without data races.
- **In-Memory Sharding:** Advanced lock contention management for extreme throughput.
- **Zero Dependencies:** Core logic relies purely on the Go standard library.

---

## 📊 Stress Test & Benchmarks

The system was stress-tested using K6 with 10 virtual users over a $20$-second period. Thanks to the in-memory architecture and optimized concurrency model, it achieved extremely low latency.

*Note: To eliminate network jitter and measure pure application processing performance, both K6 and the application were executed on the same local host.*

**Hardware Specs:** $8$ CPU Cores, $16$ GB RAM

**Results (Median of multiple test runs):**
- **Total Requests Handled:** $\sim 550,000$
- **Test Duration:** $20$ seconds
- **Throughput:** $\sim 27,500$ RPS (Requests Per Second)
- **Latency p(90):** $\sim 480 \mu s$ (microseconds)
- **Latency p(95):** $\sim 650 \mu s$ (microseconds)

---

## Architecture Highlight: Lock Contention Optimization

To achieve sub-millisecond latency ($480 \mu s$) and high throughput ($27,500$ RPS), this rate limiter utilizes an **In-Memory Sharding strategy**. 

Instead of using a single global `sync.Mutex` that would create a severe bottleneck under heavy load (Lock Contention), the internal state is partitioned into multiple concurrent-safe shards based on a hash of the User ID. This allows true parallel processing across multiple CPU cores and $O(1)$ lookup times without blocking unrelated requests.

---

## ⚠️ Known Limitations (In-Memory Version)
* **Memory Management for Inactive Users:** Currently, the in-memory implementation does not actively evict buckets for users who stop making requests. In a long-running production environment, this could lead to memory bloat. A background worker (Goroutine) could be implemented to periodically clean up expired buckets, but for true distributed scale, moving to the Redis architecture (described below) is the preferred solution.

---

## Roadmap / Future Enhancements

- [ ] **Distributed Architecture:** Implement Redis as a backend using Redis Hash Tags to ensure atomic operations and support High Availability (HA) deployments (described below).
- [ ] **Robust Error Handling & Observability:** Implement custom error types and proper error wrapping for better debugging and traceability.
- [ ] **Memory Management:** Add active eviction mechanisms (e.g., a background worker) to clean up stale buckets and prevent memory bloat over time.

---

## Future Distributed Architecture (Redis Integration)

The next phase of the roadmap is to support distributed environments (e.g., Kubernetes microservices) using **Redis** clustering. The planned architecture to use Redis focuses on maintaining atomicity and low latency:

1. **Atomic Token Deduction (Lua Scripts):** 
   Instead of using distributed locks (which introduce network round-trip overhead), the Token Bucket algorithm will be implemented entirely via **Redis Lua Scripts**. This ensures that fetching current tokens, calculating the time delta, and updating the bucket all happen in a single, atomic operation within Redis.
2. **Memory Management:** 
   User limits will be stored using Redis Hashes with a carefully calculated `TTL` (Time-To-Live). The `TTL` ensures that inactive entities are automatically evicted, preventing memory bloat (the problem In-Memory implementation had).
3. **Cluster Compatibility (Hash Tags):**
   Implementing Redis Hash Tags (e.g., user:{userID}:bucket) to ensure all keys related to a specific user map to the same Hash Slot. This is crucial for executing Lua scripts safely within a Redis Cluster environment.
4. **Fallback Strategy:** 
   Relying on infrastructure-level HA (e.g., Redis Cluster Replicas) to maintain a consistent Global State, avoiding naive in-memory fallbacks that would break the distributed limits.

---

## 🛠️ Usage (Example)
```bash
cd rate-limiter/cmd
go run main.go
curl localhost:3000/?userID=123
```