package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	envPrefix          = "zeabur"
	configDir          = ".config"
	zeaburConfigSubDir = "zeabur"
	configFile         = "cli.yaml"
)

// DefaultConfigFilePath returns the default config file path($HOME/.config/zeabur/cli.yaml)
func DefaultConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user config dir: %w", err)
	}

	return filepath.Join(homeDir, configDir, zeaburConfigSubDir, configFile), nil
}

func initViper(configPath string) {
	// create config file if not exists
	createConfigFile(configPath)

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func createConfigFile(configPath string) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configPath), 0o700); err != nil {
			panic(fmt.Errorf("could not create config directory: %w", err))
		}
		f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			panic(fmt.Errorf("could not create config file: %w", err))
		}
		if cerr := f.Close(); cerr != nil {
			panic(fmt.Errorf("could not close config file: %w", cerr))
		}
		return
	} else if err != nil {
		panic(fmt.Errorf("could not stat config file: %w", err))
	}
	// Existing files were created with os.Create (0666 & umask = 0644).
	// The YAML holds the access token; tighten it on every start.
	if err := os.Chmod(configPath, 0o600); err != nil {
		panic(fmt.Errorf("could not restrict config file mode: %w", err))
	}
}