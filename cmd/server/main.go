package main

import (
	"log"
	"os"

	"campgear/internal/httpapi"
	"campgear/internal/rental"
	"campgear/internal/storage"
)

func main() {
	path := os.Getenv("CAMP_DB_PATH")
	if path == "" {
		path = "campgear.db"
	}
	addr := os.Getenv("CAMP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	db, err := storage.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	services := rental.NewServices(db)
	server := httpapi.NewServer(services)
	log.Printf("campgear listening on %s", addr)
	if err := server.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
