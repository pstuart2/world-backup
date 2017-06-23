package cmd

import (
	"log"
	"world-backup/server/conf"
	"world-backup/server/data"

	"github.com/spf13/cobra"

	"os"
	"world-backup/server/api"

	"world-backup/server/watcher"

	"world-backup/server/fs"

	"github.com/spf13/afero"
)

var rootCmd = cobra.Command{
	Use: "example",
	Run: run,
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd.PersistentFlags().StringP("config", "c", "", "the config file to use")
	rootCmd.Flags().IntP("port", "p", 0, "the port to use")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	config, err := conf.LoadConfig(cmd)
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	logger, err := conf.ConfigureLogging(&config.LogConfig)
	if err != nil {
		log.Fatal("Failed to configure logging: " + err.Error())
	}

	aferoFs := afero.Afero{Fs: afero.NewOsFs()}

	db, dbErr := data.Open("data.json", aferoFs)
	if dbErr != nil {
		log.Fatal("Failed to open the db: " + dbErr.Error())
	}
	db.Save()

	fileSystem := fs.NewFs(aferoFs)

	w := watcher.NewWatcher(logger, config, fileSystem, db)
	w.Start()

	server := api.NewAPI(logger, config, db)
	server.SetUpRoutes()

	logger.Infof("Starting up server on port %d", config.Port)
	if err := server.Start(); err != nil {
		logger.WithError(err).Error("Error while running server")
		os.Exit(1)
	}

	logger.Info("DONE!")
}
