package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"short.ly/internal/api"
	"short.ly/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed Loading .env variables", err)
	}

	db, err := storage.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := &httpapi.Handler{
		DB: db,
	}

	log.Fatal(httpapi.StartServer(h, "8000"))
}
