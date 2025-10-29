package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AntonyKoshy/prophet-arbitrage-bot/api"
	"github.com/AntonyKoshy/prophet-arbitrage-bot/config"
)

func main() {
	// 1. Load Configuration
	// The application cannot run without its configuration, so we treat this as a fatal error.
	log.Println("Loading configuration from config.yaml...")
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("FATAL: could not load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	// 2. Initialize API Clients
	// Use the loaded configuration to create clients for the services we need.
	kalshiClient := api.NewKalshiClient(cfg.Kalshi.APIKey, cfg.Kalshi.IsDemo)
	if cfg.Kalshi.IsDemo {
		log.Println("Initialized Kalshi client in DEMO mode.")
	} else {
		log.Println("Initialized Kalshi client in PRODUCTION mode.")
	}

	// 3. Execute Application Logic
	// For this example, we'll fetch and print the events from Kalshi.
	log.Println("Fetching events from Kalshi API...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eventsResponse, err := kalshiClient.GetEvents(ctx)
	if err != nil {
		log.Fatalf("ERROR: failed to fetch Kalshi events: %v", err)
	}

	log.Printf("Successfully fetched %d events from Kalshi.", len(eventsResponse.Events))
	// Print the title of the first event as a demonstration
	if len(eventsResponse.Events) > 0 {
		fmt.Printf("→ Example Event: %s (%s)\n", eventsResponse.Events[0].Title, eventsResponse.Events[0].Ticker)
	}
}
