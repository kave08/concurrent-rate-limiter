Concurrent Rate Limiter in Go
This repository contains a simple, thread-safe rate limiter implemented in Go. It's designed to restrict the number of requests from a specific user within a defined time window, which is a common requirement for APIs and services to prevent abuse and ensure fair usage.

Features
Concurrent-Safe: Uses a sync.Mutex to ensure safe access to the rate limiter's internal state from multiple goroutines.

Configurable: The request limit and time window can be configured when initializing the rate limiter.

Efficient Cleanup: Automatically removes outdated request timestamps to prevent memory leaks over time.

How It Works
The rate limiter tracks each user's requests by storing the timestamps of their last few requests. When a new request arrives, it checks if the number of recent requests (within the defined time window) has exceeded the configured limit.

If the limit has not been reached, the request is allowed, and its timestamp is recorded.

If the limit has been reached, the request is denied.

Old timestamps are automatically cleaned up as new requests are processed, ensuring the data remains relevant to the current time window.

Usage
Prerequisites
Go 1.16 or higher

Code Example
You can use the NewRateLimiter function to create an instance and then call IsRequestAllowed to check if a request should be processed.

package main

import (
    "fmt"
    "time"
)

// The RateLimiter and UserLimit structs would be here

// NewRateLimiter function would be here

// IsRequestAllowed function would be here

func main() {
    // Create a new rate limiter that allows 3 requests per 60 seconds.
    rl := NewRateLimiter(3, 60*time.Second)
    userID := "user1"

    // Simulate requests
    for i := 0; i < 4; i++ {
        allowed := rl.IsRequestAllowed(userID, time.Now())
        fmt.Printf("Request %d: Allowed = %v\n", i+1, allowed)
        time.Sleep(1 * time.Second) // Simulate some delay between requests
    }
}

This example demonstrates how to set up the rate limiter and test its functionality. The first three requests for user1 will be allowed, while the fourth will be denied because it falls within the 60-second window and exceeds the limit of 3.

Contributing
Feel free to fork the repository, make improvements, and submit a pull request!