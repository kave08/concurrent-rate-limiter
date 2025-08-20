package main

import (
	"fmt"
	"sync"
	"time"
)

// Rate limiter function to check if a request is valid per user.
// Returns false if the limit is reached, true otherwise.
// Condition: Max 3 requests per user in the time window.
// After the window expires, the user can send requests again.
// Cleanup: Remove user from map if no requests remain in window.

type RateLimiter struct {
	users  map[string]*UserLimit
	mu     sync.Mutex
	limit  int
	window time.Duration
}

type UserLimit struct {
	requestTimes []time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		users:  make(map[string]*UserLimit),
		limit:  limit,
		window: window,
	}
}

func (r *RateLimiter) IsRequestAllowed(userID string, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		r.users[userID] = &UserLimit{requestTimes: []time.Time{now}}
		return true
	}

	cutoff := now.Add(-r.window)
	newTimes := make([]time.Time, 0, len(user.requestTimes))
	for _, t := range user.requestTimes {
		if t.After(cutoff) {
			newTimes = append(newTimes, t)
		}
	}
	user.requestTimes = newTimes

	if len(user.requestTimes) == 0 {
		delete(r.users, userID)
		r.users[userID] = &UserLimit{requestTimes: []time.Time{now}}
		return true
	}

	if len(user.requestTimes) >= r.limit {
		return false
	}

	user.requestTimes = append(user.requestTimes, now)

	return true
}

func main() {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user1"
	now := time.Now()

	for i := 0; i < 4; i++ {
		allowed := rl.IsRequestAllowed(userID, now.Add(time.Duration(i)*time.Second))
		fmt.Printf("Request %d: Allowed = %v\n", i+1, allowed)
	}
}
