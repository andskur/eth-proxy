package main

import (
	"os"

	"eth-proxy/pkg/logger"
	service "eth-proxy/proxy-service"
	"eth-proxy/proxy-service/cmd/root"
	"eth-proxy/proxy-service/cmd/serve"
)

func main() {
	app, err := service.NewApplication()
	if err != nil {
		handleError(err)
	}

	rootCmd := root.Cmd(app)
	rootCmd.AddCommand(serve.Cmd(app))

	if err := rootCmd.Execute(); err != nil {
		handleError(err)
	}
}

// handleError is
func handleError(err error) {
	logger.Log().Error(err)
	os.Exit(1)
}
