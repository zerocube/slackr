package slackr

import (
	"errors"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Allow for overrides during build time via LDFLAGS
var config_file_path string

type SlackrConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

func NewSlackrConfig() *SlackrConfig {
	var config SlackrConfig

	// Use a default config file path if one isn't provided
	if config_file_path == "" {
		config_file_path = path.Join(os.Getenv("HOME"), ".config", "slackr", "config.yaml")
	}

	// Create the config file if it doesn't exist
	if _, err := os.Stat(config_file_path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path.Dir(config_file_path), 0755)
		if err != nil {
			log.Fatalln("Unable to create parent directories:", err)
		}

		defaultConfig := SlackrConfig{}
		defaultConfigYAML, err := yaml.Marshal(defaultConfig)
		if err != nil {
			log.Fatalln("Unable to marshal default config:", err)
		}

		os.WriteFile(config_file_path, defaultConfigYAML, 0600)
	}

	// Read config from file
	file, err := os.ReadFile(config_file_path)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln(err)
	}

	return &config
}
