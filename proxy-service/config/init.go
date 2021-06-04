package config

import (
	"github.com/spf13/viper"
)

// init initialize default config params
func init() {
	// common
	viper.SetDefault("log.level", "info")

	// ethereum
	viper.SetDefault("ethereum.addr", "https://cloudflare-eth.com")
	viper.SetDefault("ethereum.wss", false)

	// grpc admin
	viper.SetDefault("grpc.host", "0.0.0.0")
	viper.SetDefault("grpc.port", 9090)
	viper.SetDefault("grpc.timeout", "120s")

	// cache
	viper.SetDefault("cachesize", 100)
}
