# Load vs Stress Testing: The Hidden Backbone of Scalable APIs

## Overview

As backend engineers, we often focus on whether our APIs return the right data, but building a reliable system goes far beyond that. What happens when hundreds or thousands of users access our API simultaneously? How does our backend handle spikes in traffic or unexpected surges?

In this article, I'll explore the critical differences between load and stress testing, explain why they are indispensable for scalable APIs, and provide a practical example using Go Fiber and k6.

## Key Differences

### Load Test (Realistic Pressure)

- Checks how the system behaves under expected traffic
- Measures response time, throughput, CPU/memory usage
- Can the API handle 1,000 concurrent users smoothly?

### Stress Test (Chaos Simulation)

- Pushes the system beyond limits to find breaking points
- Measures maximum capacity, error handling, recovery
- What happens when 10,000 users overload the API?

## Why Performance Testing Matters

The goals of performance testing:

| Goal | Description |
|------|-------------|
| Determine capacity limits | Find the maximum sustainable user load |
| Identify bottlenecks | Is the delay caused by the DB, the CPU, or the network? |
| Validate SLAs | Are 95% of API calls under 200 ms? |
| Test auto-scaling | Ensure Kubernetes or cloud autoscalers behave correctly |
| Prevent outages | Detect weak points before real users do |

## Tools for Load & Stress Testing

| Tool | Best For |
|------|----------|
| Apache JMeter | Industry standard |
| k6 | Modern load testing |
| Locust | API & Web testing |
| Artillery | REST & GraphQL APIs |
| Gatling | Enterprise scale |

## Hands-On Example: Go Fiber

Let's start with a simple Fiber server that simulates a real backend endpoint.

```go
package main

import (
	"math"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func heavyComputation() {
	// More CPU-intensive work
	for i := 0; i < 100_000_000; i++ {
		_ = math.Sin(float64(i)) * math.Sqrt(float64(i))
	}
}

func main() {
	// Use a single CPU core
	runtime.GOMAXPROCS(1)

	app := fiber.New()

	app.Get("/api/hello", func(c *fiber.Ctx) error {
		heavyComputation()
		return c.JSON(fiber.Map{
			"message": "Hello from Fiber!",
		})
	})

	app.Listen(":3000")
}
```

Now, our API is running at `http://localhost:3000/api/hello`

### Why Did We Write This Code?

Go's concurrency model is extremely powerful:

- The Go runtime runs each goroutine as a very lightweight thread
- Because of this, even with a large number of users, simple CPU or I/O operations typically execute very quickly
- Requests involving `time.Sleep(50ms)` or minor calculations are easily handled by Go's concurrency → it's hard to observe errors or timeouts

That's why we did the following:

- **Added a CPU-intensive task per request** → so that each request takes a long time to complete (`heavyComputation`)
- **Limited Go to a single core using `GOMAXPROCS(1)`** → the Go runtime cannot utilize all CPU cores, creating a CPU bottleneck

**Purpose**: Reduce the advantage of Go's concurrency and observe the system's natural limits.

In short, this code is a deliberately "heavy" simulation designed to stress the system. This allows us to measure its boundaries using stress-testing tools like k6.

## Load Testing with k6

We'll create a simple load test to simulate 10 concurrent users hitting this endpoint for 30 seconds.

### load_test.js

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,           // 10 virtual users
  duration: '30s',   // test duration
};

export default function () {
  const res = http.get('http://localhost:3000/api/hello');
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}
```

### Load Test Results

#### 1. Error Rate (http_req_failed)

```
http_req_failed................: 0.00%  0 out of 40
```

- All 40 requests succeeded
- The Fiber application did not produce any errors under load

#### 2. Response Time (http_req_duration)

```
avg=7.92s  min=6.11s  med=8.18s  max=8.68s  p(90)=8.58s  p(95)=8.64s
```

- Response times are high because `heavyComputation()` is a CPU-bound task
- Since there were no errors and all requests completed, the CPU can still handle the load

#### 3. Virtual Users (vus)

```
vus: 6   min=6   max=10   vus_max: 10
```

- The test ran with 10 virtual users, which is very low compared to a stress test
- The Fiber application responds comfortably under this load

#### 4. Iterations and Throughput

```
iterations: 40, 1.09/s
iteration_duration avg=8.92s
```

#### 5. Network

```
data_received: 5.6 kB, data_sent: 3.2 kB
```

- Only a small amount of data was sent and received
- This test is purely a CPU load test, not a network load test

## Stress Testing with k6

Now let's stress the system, gradually increasing traffic until it breaks.

### stress_test.js

```javascript
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 50 },
    { duration: '10s', target: 100 },
    { duration: '10s', target: 200 },
    { duration: '10s', target: 400 },
  ],
};

export default function () {
  const res = http.get('http://localhost:3000/api/hello');
  check(res, { 'status is 200': (r) => r.status === 200 });
}
```

### Stress Test Results

#### 1. Error Rate (http_req_failed)

```
http_req_failed................: 99.75% 19359 out of 19406
```

- 99.75% of requests failed
- Almost all requests did not receive a response from the Fiber server or timed out
- The application crashed under 400 concurrent virtual users

**Why?** Fiber running on a single core combined with the CPU-intensive `heavyComputation()` couldn't handle most requests.

#### 2. Response Times (http_req_duration)

```
avg=278.05ms  min=0s  med=0s  max=1m0s
expected_response avg=39.48s  median=40.35s  max=58.98s
p(95)=56.74s
```

- The server is CPU-blocked, queuing requests
- Under 400 VUs, the application generates severe latency

#### 3. Virtual Users (vus)

```
vus: 180  min=5  max=399
vus_max: 400
```

- The test ran with a maximum of 400 VUs
- Most of the time, 180 VUs were active while the rest were waiting or failing
- The real concurrent capacity of the application is around 180 users

## What Happens If You Don't Test?

Skipping performance tests may seem like a time-saver, but it can be disastrous:

- Without load and stress testing, we have no way of knowing the limits of our system or where bottlenecks occur
- This leaves our APIs vulnerable to latency spikes, crashes, or complete downtime when real users flood the system
- Testing proactively ensures that our backend can survive both expected and unexpected loads, safeguarding user experience and company reputation

## Real-World Consequences of Skipping Testing

### Amazon Prime Day (2018)

- **Issue**: Load tests were insufficient before launch
- **What happened**: Traffic spiked and backend servers crashed
- **Result**: 63 minutes of downtime → ~$100M loss

### Pokémon GO (2016)

- **Issue**: Niantic underestimated expected traffic (50× more users)
- **What happened**: Servers overloaded and users saw "Server is overloaded" errors for weeks
- **Result**: Massive retention loss during peak hype

### Knight Capital (2012)

- **Issue**: Lack of performance & chaos testing on trading systems
- **What happened**: A deployment bug went live and executed erroneous trades
- **Result**: $440 million loss in 45 minutes

## Conclusion: Why Your API Needs Load & Stress Testing

Load and stress testing are the hidden backbone of any scalable, production-grade API. They help uncover system weaknesses before real users experience them, ensuring your backend remains resilient, performant, and reliable under pressure.

By proactively testing your endpoints, you can:

- Identify bottlenecks
- Optimize resource usage
- Validate your SLAs
- Prepare your infrastructure for unexpected spikes in traffic

Skipping these tests may save time initially, but it risks downtime, poor user experience, and even financial loss. Incorporating regular load and stress testing into your development workflow is not just best practice—it's essential for building APIs that can grow with your users and remain trustworthy in production environments.
