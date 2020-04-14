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

	sanitizeSapControlUrl(config)

	err = validateSapControlUrl(config)
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration value")
	}

	return config, nil
}

// returns an error in case the sap-control-url config value cannot be parsed as URL
func validateSapControlUrl(config *viper.Viper) error {
	sapControlUrl := config.GetString("sap-control-url")
	if _, err := url.ParseRequestURI(sapControlUrl); err != nil {
		return errors.Wrap(err, "invalid sap-control-url value")
	}
	return nil
}

// automatically adds an http:// prefix in case it's missing from the value, to avoid the downstream consumer
// throw errors due to missing schema URL component
func sanitizeSapControlUrl(config *viper.Viper) {
	sapControlUrl := config.GetString("sap-control-url")
	hasScheme, _ := regexp.MatchString("^https?://", sapControlUrl)
	if !hasScheme {
		sapControlUrl = "http://" + sapControlUrl
		config.Set("sap-control-url", sapControlUrl)
	}
}
