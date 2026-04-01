package main

import (
	"log"
	"os"

	"github.com/stockyard-dev/stockyard-relay/internal/license"
	"github.com/stockyard-dev/stockyard-relay/internal/server"
	"github.com/stockyard-dev/stockyard-relay/internal/store"
)

func main() {
	port := getEnv("PORT", "9090")
	dataDir := getEnv("DATA_DIR", "./data")
	licenseKey := os.Getenv("RELAY_LICENSE_KEY")

	tier := "free"
	if licenseKey != "" {
		if license.Validate(licenseKey) {
			tier = "pro"
			log.Println("License valid — Pro tier active")
		} else {
			log.Println("Warning: invalid license key, running as free tier")
		}
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer db.Close()

	srv := server.New(db, tier)
	log.Printf("Stockyard Relay listening on :%s (tier: %s)", port, tier)
	log.Fatal(srv.ListenAndServe(":" + port))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
