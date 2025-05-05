package services

import (
	"sync"
	"time"
)

var (
	// token , time of expiration 
	tokenBlacklist     = make(map[string]time.Time)
	tokenBlacklistLock sync.RWMutex
)

// BlacklistToken adds a token to the blacklist
func BlacklistToken(tokenStr string) error {
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return err
	}

	tokenBlacklistLock.Lock()
	defer tokenBlacklistLock.Unlock()
	
	// add token to the blacklist with its expiration time
	tokenBlacklist[tokenStr] = claims.ExpiresAt.Time
	
	return nil
}

// IsTokenBlacklisted checks if a token is in the blacklist
func IsTokenBlacklisted(tokenStr string) bool {
	tokenBlacklistLock.RLock()
	defer tokenBlacklistLock.RUnlock()
	
	expTime, exists := tokenBlacklist[tokenStr]
	
	// cleanup expired tokens
	if exists && time.Now().After(expTime) {
		// cleanup expired token from the blacklist
		go func() {
			tokenBlacklistLock.Lock()
			delete(tokenBlacklist, tokenStr)
			tokenBlacklistLock.Unlock()
		}()
		return false
	}
	
	return exists
}

// CleanupBlacklist removes expired tokens from the blacklist
func CleanupBlacklist() {
	tokenBlacklistLock.Lock()
	defer tokenBlacklistLock.Unlock()
	
	now := time.Now()
	for token, expTime := range tokenBlacklist {
		if now.After(expTime) {
			delete(tokenBlacklist, token)
		}
	}
}