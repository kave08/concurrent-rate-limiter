package main

import (
	"testing"
	"time"
)

func TestRateLimiter_EdgeCases(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "edgeUser"
	base := time.Now()

	for i := 0; i < 3; i++ {
		if !rl.IsRequestAllowed(userID, base.Add(time.Duration(i)*time.Second)) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if rl.IsRequestAllowed(userID, base.Add(4*time.Second)) {
		t.Error("4th request within window should be denied")
	}

	if !rl.IsRequestAllowed(userID, base.Add(61*time.Second)) {
		t.Error("Request after window should be allowed")
	}

	rl.mu.Lock()
	numUsersBefore := len(rl.users)
	rl.mu.Unlock()

	rl.IsRequestAllowed(userID, base.Add(121*time.Second))

	rl.mu.Lock()
	numUsersAfter := len(rl.users)
	_, exists := rl.users[userID]
	rl.mu.Unlock()

	if !exists {
		t.Error("User should still exist after new request")
	}
	if numUsersBefore != numUsersAfter {
		t.Errorf("Number of users should remain the same after cleanup; got %d, want %d", numUsersAfter, numUsersBefore)
	}
}

func TestRateLimiter_QuickSuccession(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user1"
	startTime := time.Now()

	for i := 0; i < 4; i++ {
		requestTime := startTime.Add(time.Duration(i) * time.Second)
		allowed := rl.IsRequestAllowed(userID, requestTime)
		if i < 3 && !allowed {
			t.Errorf("Request %d at %v should be allowed", i+1, requestTime)
		}
		if i == 3 && allowed {
			t.Errorf("Request %d at %v should be denied", i+1, requestTime)
		}
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user1"
	startTime := time.Now()

	for i := 0; i < 3; i++ {
		rl.IsRequestAllowed(userID, startTime.Add(time.Duration(i)*time.Second))
	}

	requestTime := startTime.Add(61 * time.Second)
	if !rl.IsRequestAllowed(userID, requestTime) {
		t.Errorf("Request at %v should be allowed after window reset", requestTime)
	}
}

func TestRateLimiter_WindowBoundary(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user1"
	startTime := time.Now().Add(120 * time.Second)

	allowed1 := rl.IsRequestAllowed(userID, startTime)
	allowed2 := rl.IsRequestAllowed(userID, startTime.Add(10*time.Second))
	allowed3 := rl.IsRequestAllowed(userID, startTime.Add(59*time.Second))
	allowed4 := rl.IsRequestAllowed(userID, startTime.Add(61*time.Second))

	if !allowed1 || !allowed2 || !allowed3 || !allowed4 {
		t.Errorf("Expected all requests to be allowed; got %v, %v, %v, %v", allowed1, allowed2, allowed3, allowed4)
	}
}

func TestRateLimiter_MultipleUsers(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user2"
	startTime := time.Now()

	for i := 0; i < 4; i++ {
		requestTime := startTime.Add(time.Duration(i) * time.Second)
		allowed := rl.IsRequestAllowed(userID, requestTime)
		if i < 3 && !allowed {
			t.Errorf("Request %d for user2 at %v should be allowed", i+1, requestTime)
		}
		if i == 3 && allowed {
			t.Errorf("Request %d for user2 at %v should be denied", i+1, requestTime)
		}
	}
}

func TestRateLimiter_ExactWindowEnd(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user1"
	edgeStart := time.Now().Add(180 * time.Second)

	rl.IsRequestAllowed(userID, edgeStart)
	rl.IsRequestAllowed(userID, edgeStart.Add(60*time.Second))
	rl.IsRequestAllowed(userID, edgeStart.Add(60*time.Second).Add(1*time.Millisecond))
	allowedEdge := rl.IsRequestAllowed(userID, edgeStart.Add(60*time.Second).Add(2*time.Millisecond))
	if !allowedEdge {
		t.Errorf("Edge request at %v should be allowed", edgeStart.Add(60*time.Second).Add(2*time.Millisecond))
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := NewRateLimiter(3, 60*time.Second)
	userID := "user1"
	cleanupStart := time.Now().Add(240 * time.Second)

	rl.IsRequestAllowed(userID, cleanupStart) // Add one request

	rl.mu.Lock()
	numUsersBefore := len(rl.users)
	rl.mu.Unlock()

	rl.IsRequestAllowed(userID, cleanupStart.Add(61*time.Second))

	rl.mu.Lock()
	numUsersAfter := len(rl.users)
	_, exists := rl.users[userID]
	rl.mu.Unlock()

	if !exists {
		t.Error("User should still exist after cleanup and re-add")
	}
	if numUsersBefore != numUsersAfter {
		t.Errorf("Number of users should remain the same after cleanup; got %d, want %d", numUsersAfter, numUsersBefore)
	}
}