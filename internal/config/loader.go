package config

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadConfigurationAutomatically(config *viper.Viper) error {
	config.SetConfigName("sap_host_exporter")
	config.AddConfigPath("./")
	config.AddConfigPath("$HOME/.config/")
	config.AddConfigPath("/etc/")
	config.AddConfigPath("/usr/etc/")

	err := config.ReadInConfig()
	if err == nil {
		log.Info("Using config file: ", config.ConfigFileUsed())
		return nil
	}

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Infof("Could not discover configuration file: %s", err)
		log.Info("Default configuration values will be used")
		return nil
	}

	return errors.Wrap(err, "could not load automatically discovered config file")
}

// loads configuration from an explicit file path
func loadConfigurationFromFile(config *viper.Viper, path string) error {
	// we hard-code the config type to yaml, otherwise ReadConfig will not load the values
	// see https://github.com/spf13/viper/issues/316
	config.SetConfigType("yaml")

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not open file")
	}
	defer file.Close()

	err = config.ReadConfig(file)
	if err != nil {
		return errors.Wrap(err, "could not read file")
	}

	log.Info("Using config file: ", path)

	return nil
}
