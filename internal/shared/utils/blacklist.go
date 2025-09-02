package utils

import "sync"

// Blacklist stores invalidated JWT tokens
type Blacklist struct {
	tokens map[string]struct{}
	mu     sync.RWMutex
}

var blacklist = &Blacklist{
	tokens: make(map[string]struct{}),
}

// Add adds a token to the blacklist
func Add(token string) {
	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()
	blacklist.tokens[token] = struct{}{}
}

// IsBlacklisted checks if a token is blacklisted
func IsBlacklisted(token string) bool {
	blacklist.mu.RLock()
	defer blacklist.mu.RUnlock()
	_, exists := blacklist.tokens[token]
	return exists
}