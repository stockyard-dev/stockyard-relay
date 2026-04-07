package main

import (
	"flag"
	"fmt"
	"github.com/stockyard-dev/stockyard-relay/internal/server"
	"github.com/stockyard-dev/stockyard-relay/internal/store"
	"log"
	"net/http"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "")
	dataFlag := flag.String("data", "", "")
	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8620"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./relay-data"
	}
	if *portFlag != "" {
		port = *portFlag
	}
	if *dataFlag != "" {
		dataDir = *dataFlag
	}
	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("relay: %v", err)
	}
	defer db.Close()
	srv := server.New(db, server.DefaultLimits(), dataDir)
	fmt.Printf("\n  Relay — Self-hosted webhook relay\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Receive:    http://localhost:%s/hook/{channel}\n  Data:       %s\n  ─────────────────────────────────\n  Questions? hello@stockyard.dev\n\n", port, port, port, dataDir)
	log.Printf("relay: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
