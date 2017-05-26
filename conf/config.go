package conf

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config the application's configuration
type Config struct {
	Port      int64         `json:"port"`
	WatchDir  string        `json:"watchDir"`
	BackupDir string        `json:"backupDir"`
	Path      string        `json:"path"`
	LogConfig LoggingConfig `json:"log"`
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	viper.SetConfigType("json") // necessary in the event that we are using file-less config

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./")
	}

	if err := viper.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return populateConfig(config)
}
