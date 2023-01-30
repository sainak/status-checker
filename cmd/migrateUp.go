package cmd

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sainak/status-checker/core/config"
	"github.com/sainak/status-checker/core/logger"
	"github.com/spf13/cobra"
)

var migrateUpCmd *cobra.Command

func init() {
	migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "migrate to v1 command",
		Long:  `Command to install version 1 of our application`,
		Run: func(cmd *cobra.Command, args []string) {
			config.GetConfig()

			logger.Info("Running migrate up command")

			db, _ := sql.Open("postgres", config.GetDBurl())
			defer db.Close()

			dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
			if err != nil {
				logger.Error(fmt.Sprintf("instance error: %v \n", err))
			}

			fileSource, err := (&file.File{}).Open("file://migrations")
			if err != nil {
				logger.Error(fmt.Sprintf("opening file error: %v \n", err))
			}

			m, err := migrate.NewWithInstance("file", fileSource, "postgres", dbDriver)
			if err != nil {
				logger.Error(fmt.Sprintf("migrate error: %v \n", err))
			}

			if err = m.Up(); err != nil {
				logger.Error(fmt.Sprintf("migrate up error: %v \n", err))
			}

			logger.Info("Migrate up done with success")
		},
	}

	migrateCmd.AddCommand(migrateUpCmd)
}
