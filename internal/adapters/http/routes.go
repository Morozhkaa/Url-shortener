package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	Status int
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		wrapper := &responseWriter{ResponseWriter: rw}
		start := time.Now()
		h.ServeHTTP(wrapper, r)
		elapsed := time.Since(start)

		log.Printf("%s %s: %d %s", r.Method, r.URL, wrapper.Status, elapsed)
	})
}

func initRouter(a *Adapter, r *mux.Router) {
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"*"}),
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	r.Use(loggingMiddleware)
	r.Use(cors)
	r.HandleFunc("/{shortUrl:\\w{5}}", a.ResolveURL).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", a.CreateShortcut).Methods(http.MethodPost)
}
