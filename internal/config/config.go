package config

import (
	"github.com/spf13/viper"
)

const (
	DefaultBaseURL = "https://d2v1fiwobg9w6.cloudfront.net"
	DefaultInfoURL = "https://api.hyperliquid.xyz/info"
)

type Config struct {
	BaseURL string
	InfoURL string
	Format  string
}

func New() *Config {
	return &Config{
		BaseURL: viper.GetString("base_url"),
		InfoURL: viper.GetString("info_url"),
	}
}
