package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"short.ly/db"
)

type Handler struct {
	pool      *pgxpool.Pool
	cache     *redis.Client
	cache_ttl time.Duration
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type RedirectResponse struct {
	URL string `json:"url"`
}

func NewHandler(p *pgxpool.Pool, c *redis.Client) *Handler {
	c_ttl, err := db.CacheTTL()
	if err != nil {
		log.Println(err)
	}
	return &Handler{
		pool:      p,
		cache:     c,
		cache_ttl: c_ttl,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if short_url, err := db.CreateLink(h.pool, r.Context(), req.URL); err != nil {
		fmt.Println(r.Context().Value(http.LocalAddrContextKey))
		log.Println("create link failed:", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ShortenResponse{ShortURL: short_url})
		return
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	redirect_url, err := h.cache.Get(r.Context(), code).Result()
	if err == nil {
		log.Println("cache hit for", code)
	} else {
		if errors.Is(err, redis.Nil) {
			log.Println("cache miss for", code)
		} else {
			log.Println("redis error, falling back to db:", err)
		}

		redirect_url, err = db.GetLongURL(h.pool, r.Context(), code)
		if errors.Is(err, db.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			log.Println("Error fetching url", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		if err := h.cache.Set(r.Context(), code, redirect_url, h.cache_ttl).Err(); err != nil {
			log.Println("cache write failed:", err)
		}
	}

	http.Redirect(w, r, redirect_url, http.StatusMovedPermanently)
}
