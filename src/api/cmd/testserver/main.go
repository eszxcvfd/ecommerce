// Command testserver starts the catalog API with seeded data for e2e tests.
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
	repo := catalog.NewMemoryRepo(seed())
	log.Printf("Test API server listening on :%s", port)
	if err := catalog.StartServer(":"+port, repo); err != nil {
		log.Fatal(err)
	}
}

func seed() []catalog.SanPhamSo { return catalog.SeedData() }
