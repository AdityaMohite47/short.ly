package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"short.ly/db"
)

type Handler struct {
	pool *pgxpool.Pool
	cache *redis.Client
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

func NewHandler(p *pgxpool.Pool , c *redis.Client) *Handler {
	return &Handler{
		pool: p,
		cache: c,
	}	
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if short_url, err := db.CreateLink(h.pool , r.Context(), req.URL); err != nil {
		fmt.Println(r.Context().Value(http.LocalAddrContextKey))
		log.Println("create link failed:", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ShortenResponse{ShortURL: short_url})
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
    
	redirect_url, err := db.GetLongURL(h.pool, r.Context(), code)
	if errors.Is(err, db.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Println("Error fetching url", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect_url, http.StatusMovedPermanently)
}
