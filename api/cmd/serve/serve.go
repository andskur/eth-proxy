package serve

import (
	"fmt"

	"github.com/spf13/cobra"

	"eth-proxy/api"
)

// Cmd run application by calling Serve method
func Cmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run Application",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Init(); err != nil {
				return err
			}

			return app.Serve()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf(app.Version())
		},
	}
}
