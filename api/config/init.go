package config

import (
	"github.com/spf13/viper"
)

// init initialize default config params
func init() {
	// common
	viper.SetDefault("log.level", "info")

	// grpc admin
	viper.SetDefault("api.host", "proxy-service")
	viper.SetDefault("api.port", 9090)
	viper.SetDefault("api.timeout", "60s")

	// http
	viper.SetDefault("http.host", "0.0.0.0")
	viper.SetDefault("http.port", 8888)
	viper.SetDefault("http.timeout", "60s")
}
