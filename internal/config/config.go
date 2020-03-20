package config

import (
	"net/url"
	"regexp"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func New() (*viper.Viper, error) {
	config := viper.New()

	err := config.BindPFlags(flag.CommandLine)
	if err != nil {
		return nil, errors.Wrap(err, "could not bind config to CLI flags")
	}

	// try to get the "config" value from the bound "config" CLI flag
	path := config.GetString("config")
	if path != "" {
		// try to manually load the configuration from the given path
		err = loadConfigurationFromFile(config, path)
	} else {
		// otherwise try viper's auto-discovery
		err = loadConfigurationAutomatically(config)
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not load configuration file")
	}

	setLogLevel(config.GetString("log-level"))

	err = validateConfigValues(config)
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration value")
	}

	sanitizeConfigValues(config)

	return config, nil
}

func validateConfigValues(config *viper.Viper) error {
	sapControlUrl := config.GetString("sap-control-url")
	if _, err := url.ParseRequestURI(sapControlUrl); err != nil {
		return errors.Wrap(err, "invalid config value for sap-control-url")
	}
	return nil
}

func sanitizeConfigValues(config *viper.Viper) {
	sapControlUrl := config.GetString("sap-control-url")
	hasScheme, _ := regexp.MatchString("^https?://", sapControlUrl)
	if !hasScheme {
		sapControlUrl = "http://" + sapControlUrl
		config.Set("sap-control-url", sapControlUrl)
	}
}
