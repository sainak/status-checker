package cmd

import (
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database to schema version",
	Long:  `migrate cmd is used for database migration: migrate < up | down >`,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
