package services

import (
	"sync"
	"time"
)

var (
	// Черный список токенов: ключ - токен, значение - время истечения
	tokenBlacklist     = make(map[string]time.Time)
	tokenBlacklistLock sync.RWMutex
)

// BlacklistToken добавляет токен в черный список
func BlacklistToken(tokenStr string) error {
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return err
	}

	tokenBlacklistLock.Lock()
	defer tokenBlacklistLock.Unlock()
	
	// Добавляем токен в черный список до времени его истечения
	tokenBlacklist[tokenStr] = claims.ExpiresAt.Time
	
	return nil
}

// IsTokenBlacklisted проверяет, находится ли токен в черном списке
func IsTokenBlacklisted(tokenStr string) bool {
	tokenBlacklistLock.RLock()
	defer tokenBlacklistLock.RUnlock()
	
	expTime, exists := tokenBlacklist[tokenStr]
	
	// Очистка устаревших токенов
	if exists && time.Now().After(expTime) {
		// Отложенная очистка (в реальном приложении можно сделать отдельную горутину)
		go func() {
			tokenBlacklistLock.Lock()
			delete(tokenBlacklist, tokenStr)
			tokenBlacklistLock.Unlock()
		}()
		return false
	}
	
	return exists
}

// CleanupBlacklist удаляет просроченные токены из черного списка
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