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

var migrateDownCmd *cobra.Command

func init() {
	migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "migrate from v2 to v1",
		Long:  `Command to downgrade database from v2 to v1`,
		Run: func(cmd *cobra.Command, args []string) {
			config.GetConfig()

			logger.Info("Running migrate down command")

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

			m, err := migrate.NewWithInstance("file", fileSource, config.GetDBurl(), dbDriver)
			if err != nil {
				logger.Error(fmt.Sprintf("migrate error: %v \n", err))
			}

			if err = m.Down(); err != nil {
				logger.Error(fmt.Sprintf("migrate down error: %v \n", err))
			}

			logger.Info("Migrate down done with success")
		},
	}

	migrateCmd.AddCommand(migrateDownCmd)
}
