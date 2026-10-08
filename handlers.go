package main

import (
	"encoding/json"
	// "errors"
	"log"
	"net/http"
)

type Handler struct {
	DB      *DBconn
	// BaseURL string
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if short_url, err := h.DB.getShortURL(r.Context(), req.URL); err != nil {
		log.Println("create link failed:", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ShortenResponse{ShortURL:short_url})
	}

}