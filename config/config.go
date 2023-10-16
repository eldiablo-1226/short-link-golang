package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	DatabaseUrl string
	Domain      string
	Username    string
	Password    string
}

var AppConfig *Config

func Init() error {
	configPath, err := ParseFlags()

	if err != nil {
		return err
	}

	var conf Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		log.Fatalf("configuration not found %v", configPath)
		return err
	}

	AppConfig = &conf

	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.toml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
