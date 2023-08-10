package inmemoryimpl

import (
	"context"
	"log"
	"sync"
	"url-shortener/internal/domain/models"
	"url-shortener/internal/domain/usecases"
)

func NewManager() *Manager {
	return &Manager{
		urlShortcuts: make(map[string]string),
	}
}

type Manager struct {
	mu           sync.RWMutex
	urlShortcuts map[string]string // [short url key] -> full url
}

func (m *Manager) insertIfKeyIsNotUsed(key string, fullURL string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.urlShortcuts[key]; ok {
		return false
	}
	m.urlShortcuts[key] = fullURL
	return true
}

func (m *Manager) CreateShortcut(ctx context.Context, fullURL string) (string, error) {
	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		key := usecases.GenerateKey()
		succeeded := m.insertIfKeyIsNotUsed(key, fullURL)
		if !succeeded {
			log.Printf("Got collision for key %s. Retry key generation...", key)
			continue
		}
		return key, nil
	}
	return "", models.ErrKeyGenerationFailed
}

func (m *Manager) ResolveShortcut(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, found := m.urlShortcuts[key]
	if !found {
		return "", models.ErrNotFound
	}
	return url, nil
}
