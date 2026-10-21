package main

import (
	"context"
	"log"
	"os"
	"github.com/joho/godotenv"
	"short.ly/api"
	"short.ly/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed Loading .env variables", err)
	}


	dbpool, err := db.Connectpsql(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	rdb , err := db.Connectredis(context.Background() , os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

	log.Fatal(api.StartServer(api.NewHandler(dbpool , rdb), "8000"))
}
