package main

import (
	"log"
	"os"

	"ecommerce/api/catalog"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	repo := catalog.NewMemoryRepo(nil) // production would wire a real store
	log.Printf("API server listening on :%s", port)
	if err := catalog.StartServer(":"+port, repo); err != nil {
		log.Fatal(err)
	}
}
