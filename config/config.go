package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration struct for the application.
type Config struct {
	Kalshi     KalshiConfig     `yaml:"kalshi"`
	Polymarket PolymarketConfig `yaml:"polymarket"`
	Bot        BotConfig        `yaml:"bot"`
}

// KalshiConfig holds the configuration specific to the Kalshi API.
type KalshiConfig struct {
	APIKey string `yaml:"api_key"`
	IsDemo bool   `yaml:"is_demo"`
}

// PolymarketConfig holds the configuration specific to the Polymarket API.
type PolymarketConfig struct {
	// Example: You might need a private key for a wallet to trade
	PrivateKey string `yaml:"private_key"`
}

// BotConfig holds general configuration for the arbitrage bot's logic.
type BotConfig struct {
	MinProfitThreshold float64 `yaml:"min_profit_threshold"`
	MaxTradeValue      float64 `yaml:"max_trade_value"`
}

// LoadConfig reads configuration from a YAML file at the given path.
func LoadConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &cfg, nil
}
