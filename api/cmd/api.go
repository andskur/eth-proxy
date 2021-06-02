package main

import (
	"os"

	"eth-proxy/api"
	"eth-proxy/api/cmd/root"
	"eth-proxy/api/cmd/serve"
	"eth-proxy/pkg/logger"
)

func main() {
	app, err := api.NewApplication()
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
