package serve

import (
	"fmt"

	"github.com/spf13/cobra"

	service "eth-proxy/proxy-service"
)

// Cmd run application by calling Serve method
func Cmd(app *service.App) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run Application",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Init(); err != nil {
				return fmt.Errorf("application initialisation: %w", err)
			}

			return app.Serve()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf(app.Version())
		},
	}
}
