package main

import (
	"log"
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.Shorten)
	// mux.HandleFunc("GET /{code}", h.Redirect)
	return mux
}

func StartServer(h *Handler, port string) error {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      h.Routes(),
	}

	log.Println("listening on port", port)
	return server.ListenAndServe()
}