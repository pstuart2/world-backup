package conf

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config the application's configuration
type Config struct {
	Port          int64         `json:"port"`
	WatchDirs     []string      `json:"watchDirs"`
	BackupDir     string        `json:"backupDir"`
	CheckInterval string        `json:"checkInterval"`
	LogConfig     LoggingConfig `json:"log"`
	StaticRoot    string        `json:"staticRoot"`
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	viper.SetConfigType("json")

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

	return cleanConfig(config), nil
}

func cleanConfig(c *Config) *Config {
	c.BackupDir = os.ExpandEnv(c.BackupDir)

	var paths []string
	for i := range c.WatchDirs {
		paths = append(paths, os.ExpandEnv(c.WatchDirs[i]))
	}

	c.WatchDirs = paths

	return c
}
