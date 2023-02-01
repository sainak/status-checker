package cmd

import (
	"github.com/sainak/status-checker/app"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run the server",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
