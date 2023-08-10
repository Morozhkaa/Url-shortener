package application

import (
	"context"
	"fmt"
	"url-shortener/internal/adapters/http"
	"url-shortener/internal/domain/models"
	"url-shortener/internal/domain/usecases/inmemoryimpl"
	"url-shortener/internal/domain/usecases/mongoimpl"
	rediscached "url-shortener/internal/domain/usecases/rediscachedimpl"
	"url-shortener/internal/ports"
	"url-shortener/internal/ratelimit"

	"github.com/go-redis/redis/v8"
)

type App struct {
	opts          AppOptions
	shutdownFuncs []func(ctx context.Context) error
}

type AppOptions struct {
	Mode      models.Mode
	HTTP_port int
	MongoAddr string
	RedisAddr string
}

func New(opts AppOptions) *App {
	return &App{
		opts: opts,
	}
}

func (app *App) Start() error {
	var manager ports.Manager
	redisClient := redis.NewClient(&redis.Options{Addr: app.opts.RedisAddr})
	limiterFactory := ratelimit.NewFactory(redisClient)

	switch app.opts.Mode {
	case models.ModeMongo:
		manager = mongoimpl.NewManager(app.opts.MongoAddr)
	case models.ModeCached:
		mongoManager := mongoimpl.NewManager(app.opts.MongoAddr)
		manager = rediscached.NewManager(redisClient, mongoManager)
	default:
		manager = inmemoryimpl.NewManager()
	}

	optsAdapter := http.AdapterOptions{HTTP_port: app.opts.HTTP_port}
	s, err := http.New(manager, optsAdapter, limiterFactory)
	if err != nil {
		return fmt.Errorf("server not started: %w", err)
	}

	app.shutdownFuncs = append(app.shutdownFuncs, s.Stop)
	err = s.Start()
	if err != nil {
		return fmt.Errorf("server not started: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	var err error
	for i := len(a.shutdownFuncs) - 1; i >= 0; i-- {
		err = a.shutdownFuncs[i](ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
