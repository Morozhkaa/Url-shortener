package ports

import "context"

type Manager interface {
	CreateShortcut(ctx context.Context, fullURL string) (string, error)
	ResolveShortcut(ctx context.Context, key string) (string, error)
}
