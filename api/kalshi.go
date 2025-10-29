package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// kalshiAPIBaseURLV2 is the base URL for the production Kalshi Trading API.
	kalshiAPIBaseURLV2 = "https://trading-api.kalshi.com/v2"
	// kalshiAPIDemoBaseURLV2 is the base URL for the demo Kalshi Trading API.
	// Note: The demo environment uses API v1.
	kalshiAPIDemoBaseURLV2 = "https://demo-api.kalshi.co/trade-api/v2"
)

// KalshiClient is a client for interacting with the Kalshi API.
type KalshiClient struct {
	BaseURL    string
	HTTPClient *http.Client
	apiKey     string
}

// NewKalshiClient creates a new client for the Kalshi API.
// Pass an empty string for apiKey for unauthenticated requests.
// Set isDemo to true to use the demo environment.
func NewKalshiClient(apiKey string, isDemo bool) *KalshiClient {
	baseURL := kalshiAPIBaseURLV2
	if isDemo {
		baseURL = kalshiAPIDemoBaseURLV2
	}
	return &KalshiClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// --- Response Structs ---
// These structs are used to unmarshal the JSON responses from Kalshi.
// It's good practice to define them based on the API documentation.

type GetEventsResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	Ticker            string   `json:"ticker"`
	SeriesTicker      string   `json:"series_ticker"`
	Title             string   `json:"title"`
	Status            string   `json:"status"`
	Markets           []Market `json:"markets"`
	ExpirationTS      int64    `json:"expiration_ts"`
	SettlementSources []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"settlement_sources"`
}

type Market struct {
	Ticker        string `json:"ticker"`
	Title         string `json:"title"`
	YesBid        int    `json:"yes_bid"`
	YesAsk        int    `json:"yes_ask"`
	NoBid         int    `json:"no_bid"`
	NoAsk         int    `json:"no_ask"`
	LastPrice     int    `json:"last_price"`
	PreviousPrice int    `json:"previous_price"`
}

// --- API Methods ---

// GetEvents fetches a list of events from Kalshi.
// This corresponds to the GET /events endpoint.
func (c *KalshiClient) GetEvents(ctx context.Context) (*GetEventsResponse, error) {
	// 1. Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/events", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvents request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Api-Key "+c.apiKey)
	}

	// 2. Execute the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetEvents request: %w", err)
	}
	defer res.Body.Close()

	// 3. Check for successful status code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetEvents request failed with status: %s", res.Status)
	}

	// 4. Decode the JSON response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var eventsResponse GetEventsResponse
	if err := json.Unmarshal(body, &eventsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events response: %w", err)
	}

	return &eventsResponse, nil
}
