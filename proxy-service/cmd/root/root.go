package root

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"eth-proxy/pkg/logger"
	service "eth-proxy/proxy-service"
	"eth-proxy/proxy-service/config"
)

// config path
var configPath string

// Cmd is root level command which parsing config file and env vars
func Cmd(app *service.App) *cobra.Command {
	cmd := &cobra.Command{
		Short:   "eth proxy service",
		Version: app.Version(),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return parseCfg(app.Config())
		},
		TraverseChildren: true,
	}

	cmd.PersistentFlags().StringVarP(
		&configPath, "config", "c", "", "config file path",
	)
	cmd.SetVersionTemplate(app.Version())

	return cmd
}

// parseCfg parse config params inti structure Scheme
func parseCfg(cfg *config.Scheme) error {
	// set temporary logger
	log := logger.Log()

	// Trying to open config file
	if configPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configPath)

		// Attempts to load configuration
		if err := viper.ReadInConfig(); err != nil {
			log.Warning(err)
		}
	}

	// set config via env vars
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return viper.Unmarshal(cfg)
}
