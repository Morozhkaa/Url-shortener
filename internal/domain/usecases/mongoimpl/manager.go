package mongoimpl

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/internal/domain/models"
	"url-shortener/internal/domain/usecases"
	"url-shortener/internal/ports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "shortUrls"
const collName = "urls"

type manager struct {
	client *mongo.Client
	urls   *mongo.Collection
}

func NewManager(mongoURL string) *manager {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}
	collection := client.Database(dbName).Collection(collName)

	return &manager{
		client: client,
		urls:   collection,
	}
}

var _ ports.Manager = (*manager)(nil)

func (s *manager) CreateShortcut(ctx context.Context, url string) (string, error) {
	const attemptsCount = 5
	for attempt := 0; attempt < attemptsCount; attempt++ {
		key := usecases.GenerateKey()
		item := models.UrlItem{
			Key: key,
			URL: url,
		}

		_, err := s.urls.InsertOne(ctx, item)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}
			return "", fmt.Errorf("something went wrong: %v - %w", err, models.ErrStorage)
		}
		return key, nil
	}
	return "", fmt.Errorf("too much attempts during inserting - %w", models.ErrStorage)
}

func (s *manager) ResolveShortcut(ctx context.Context, key string) (string, error) {
	var result models.UrlItem
	err := s.urls.FindOne(ctx, bson.M{"_id": key}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("no document with key %v - %w", key, models.ErrNotFound)
		}
		return "", fmt.Errorf("somehting went wrong - %w", models.ErrStorage)
	}
	return result.URL, nil
}
