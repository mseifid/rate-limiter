# High-Performance Rate Limiter in Go

An ultra fast, concurrent-safe rate limiting service written in Go, implementing the **Token Bucket** algorithm. Current implementation is in-memory and for a single-node, however there is a roadmap to be extended for distributed scaling.

## Features

- **Token Bucket Algorithm:** This algorithm guarantees smooth and predictable rate limiting.
- **Concurrent-Safe:** Built and optimized from scratch to handle massive concurrent requests without data races.
- **In-Memory Sharding:** Handled lock contention primarily using **sharding** for extreme throughput.
- **Zero Dependencies:** Core logic relies purely on the Go standard library.

---

## 📊 Stress Test & Benchmarks

The system was stress-tested using K6 with 10 virtual users over a $20$-second period. As I have used in-memory architecture and optimized concurrency model, it achieved extremely low latency.

*Note: To eliminate network latency and measure pure application processing performance, both K6 and the application were executed on the same local host.*

**Hardware Specs:** $8$ CPU Cores, $16$ GB RAM

**Results (Median of multiple test runs):**
- **Total Requests Handled:** $\sim 550,000$
- **Test Duration:** $20$ seconds
- **Throughput:** $\sim 27,500$ RPS (Requests Per Second)
- **Latency p(90):** $\sim 480 \mu s$ (microseconds)
- **Latency p(95):** $\sim 650 \mu s$ (microseconds)

---

## Architecture Highlight: Lock Contention Optimization

To achieve this microseconds latency ($480 \mu s$) and high throughput ($27,500$ RPS), an **In-Memory Sharding strategy** is used. 

Instead of using a single global `sync.Mutex` that creates a severe bottleneck under heavy load (Lock Contention), the user's buckets are distributed into multiple concurrent-safe shards based on the hash of the User ID. This architecture makes the lookup time of $O(1)$ without blocking the unrelated users (requests).

---

## ⚠️ Known Limitations (In-Memory Version)
* **Memory Management for Inactive Users:** Currently, the in-memory implementation does not actively remove buckets for users who stop making requests for a long time. In a long-running production environment, this could lead to memory bloat. The cleanest approach is to have a `TTL` (Time to live) within every bucket, letting the bucket be removed when the `TTL` expires. This needs a more sophisticated memory (like Redis).
Another improvement could be the re-creation of the bucket when the request arrives and bucket layer sees that a long time has passed from the last refill of the bucket, so it would not refill the bucket, just deletes and re-creates it.

---

## Roadmap / Future Enhancements

- [ ] **Distributed Architecture:** Implement Redis logic using it's **Hash Tags** to support distributed concurrent-safe flow (described below).
- [ ] **Robust Error Handling:** Implement custom error types and proper error wrapping for better debugging and traceability.
- [ ] **Memory Management:** Add bucket re-creation logic to avoid redundant calculations.

---

## Future Distributed Architecture (Redis Integration)

The roadmap's major step is to support distributed environments using **Redis** clustering. Atomicity is the primary focus of the proposed architecture:

1. **Atomic Token Deduction (Lua Scripts):** 
   The Token Bucket algorithm will be implemented entirely using **Redis Lua Scripts**. This ensures that all bucket logic stages and calculations happen in a single, atomic operation within Redis.
2. **Memory Management:** 
   User's buckets will be stored in Redis using Redis Hashes (actually Hash Tags) with a calculated `TTL`. The `TTL` ensures that inactive buckets will be automatically removed, preventing memory bloat (the problem In-Memory implementation had).
3. **Cluster Compatibility (Hash Tags):**
   Redis implementation will be using Hash Tags (e.g., user:{userID}:bucket) to ensure all keys related to a specific user, map to the same Hash Slot. This is crucial for executing Lua scripts safely within a Redis cluster environment.

---

## 🛠️ Usage (Example)
```bash
cd rate-limiter/cmd
go run main.go
curl localhost:3000/?userID=123
```