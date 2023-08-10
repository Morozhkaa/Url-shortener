package rediscached

import (
	"context"
	"fmt"
	"log"
	"time"
	"url-shortener/internal/domain/models"
	"url-shortener/internal/ports"

	"github.com/go-redis/redis/v8"
)

type manager struct {
	client            *redis.Client
	persistentManager ports.Manager
}

func NewManager(client *redis.Client, persistentManager ports.Manager) *manager {
	return &manager{
		client:            client,
		persistentManager: persistentManager,
	}
}

var _ ports.Manager = (*manager)(nil)

func (m *manager) CreateShortcut(ctx context.Context, fullURL string) (string, error) {
	key, err := m.persistentManager.CreateShortcut(ctx, fullURL)
	if err != nil {
		return "", err
	}
	m.store(ctx, key, fullURL)
	return key, nil
}

func (m *manager) ResolveShortcut(ctx context.Context, key string) (string, error) {
	result := m.client.Get(ctx, m.redisKey(key))
	switch fullURL, err := result.Result(); {
	case err == redis.Nil:
	// continue execution
	case err != nil:
		return "", fmt.Errorf("%w: failed to get value from redis due to error %s", models.ErrStorage, err)
	default:
		log.Printf("Successfully obtained url from cache for key %s", key)
		return fullURL, nil
	}

	log.Printf("Loading url by key %s from persistent storage", key)
	fullURL, err := m.persistentManager.ResolveShortcut(ctx, key)
	if err != nil {
		return "", err
	}
	m.store(ctx, key, fullURL)
	return fullURL, nil
}

func (m *manager) store(ctx context.Context, shortKey string, fullURL string) {
	if err := m.client.Set(ctx, m.redisKey(shortKey), fullURL, time.Hour).Err(); err != nil {
		log.Printf("Failed to insert key %s into cache due to an error: %+v\n", shortKey, err)
	}
}

func (m *manager) redisKey(shortKey string) string {
	return "surl:" + shortKey // add a prefix not to collide with other data stored in the same redis
}
