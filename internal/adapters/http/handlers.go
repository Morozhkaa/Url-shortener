package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/domain/models"
	"url-shortener/internal/ratelimit"
)

type CreateShortcutRequest struct {
	Url string `json:"url"`
}

type CreateShortcutResponse struct {
	Key string `json:"key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func isRateLimited(limiter *ratelimit.Limiter, rw http.ResponseWriter, r *http.Request) bool {
	canDo, err := limiter.CanDoAt(r.Context(), time.Now())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return true
	}
	if !canDo {
		http.Error(rw, "rate limit exceeded", http.StatusTooManyRequests)
		return true
	}
	return false
}

func (a *Adapter) CreateShortcut(rw http.ResponseWriter, r *http.Request) {
	if isRateLimited(a.createLimit, rw, r) {
		return
	}
	var data CreateShortcutRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := a.manager.CreateShortcut(r.Context(), data.Url)
	var status int
	var response interface{}
	if err != nil {
		status = http.StatusInternalServerError
		response = ErrorResponse{
			Error: err.Error(),
		}
	} else {
		status = http.StatusOK
		response = CreateShortcutResponse{
			Key: key,
		}
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(rawResponse)
}

func (a *Adapter) ResolveURL(rw http.ResponseWriter, r *http.Request) {
	if isRateLimited(a.resolveLimit, rw, r) {
		return
	}
	key := strings.Trim(r.URL.Path, "/")

	url, err := a.manager.ResolveShortcut(r.Context(), key)
	if errors.Is(err, models.ErrNotFound) {
		http.NotFound(rw, r)
	}
	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusPermanentRedirect)
}
