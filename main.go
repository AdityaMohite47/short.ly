package main

import (
	"log"
)

func main() {
	db := &DBconn{}
	db.ConnectDB()

	h := &Handler{
		DB:   db,
	}

	log.Fatal(StartServer(h, "8000"))
}