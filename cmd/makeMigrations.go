package cmd

import (
	"fmt"
	"github.com/sainak/status-checker/core/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	migrationPath = "./migrations"
)

var createMigrationFileCmd *cobra.Command

func init() {
	createMigrationFileCmd = &cobra.Command{
		Use:   "makemigrations",
		Short: "create migration files",
		Run: func(cmd *cobra.Command, args []string) {
			var filename string
			if len(args) == 0 {
				logger.Fatal("filename not provided")
			} else {
				filename = args[0]
			}

			var lastMigrationNum int
			re := regexp.MustCompile("^([0-9]+)_.*") // to capture the number from number_name.sql
			err := filepath.Walk(migrationPath, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				match := re.FindStringSubmatch(info.Name())
				if len(match) > 0 {
					if n, err := strconv.Atoi(match[1]); err == nil && n > lastMigrationNum {
						lastMigrationNum = n
					}
				}
				return nil
			})
			if err != nil {
				logger.Fatal(err)
			}

			lastMigrationNum++
			upMigrationFilePath := fmt.Sprintf("%s/%d_%s.up.sql", migrationPath, lastMigrationNum, filename)
			downMigrationFilePath := fmt.Sprintf("%s/%d_%s.down.sql", migrationPath, lastMigrationNum, filename)

			err = createFile(upMigrationFilePath)
			if err != nil {

				return
			}
			err = createFile(downMigrationFilePath)
			if err != nil {
				err = os.Remove(upMigrationFilePath)
				return
			}

			logger.WithFields(logrus.Fields{
				"up":   upMigrationFilePath,
				"down": downMigrationFilePath,
			}).Info("Created migration files")

			return
		},
	}

	rootCmd.AddCommand(createMigrationFileCmd)
}

func createFile(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	err = f.Close()
	return
}
