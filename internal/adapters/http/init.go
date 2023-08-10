package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
	"url-shortener/internal/ports"
	"url-shortener/internal/ratelimit"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

type Adapter struct {
	s       *http.Server
	l       net.Listener
	manager ports.Manager

	createLimit  *ratelimit.Limiter
	resolveLimit *ratelimit.Limiter
}

type AdapterOptions struct {
	HTTP_port int
}

func New(manager ports.Manager, opts AdapterOptions, limiterFactory *ratelimit.Factory) (*Adapter, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.HTTP_port))
	if err != nil {
		return nil, fmt.Errorf("server start failed: %w", err)
	}
	router := mux.NewRouter()
	server := http.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a := Adapter{
		s:            &server,
		l:            l,
		manager:      manager,
		createLimit:  limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
		resolveLimit: limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
	}
	initRouter(&a, router)
	return &a, nil
}

func (a *Adapter) Start() error {
	eg := &errgroup.Group{}
	eg.Go(func() error {
		return a.s.Serve(a.l)
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		err = a.s.Shutdown(ctx)
	})
	return err
}
