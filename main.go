package main

import (
	"fmt"
	"net/http"

	"github.com/hooklift/gowsdl/soap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/collector/enqueue_server"
	"github.com/SUSE/sap_host_exporter/collector/start_service"
	"github.com/SUSE/sap_host_exporter/internal"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func init() {
	config.SetConfigName("sap_host_exporter")
	config.AddConfigPath("./")
	config.AddConfigPath("$HOME/.config/")
	config.AddConfigPath("/etc/")
	config.AddConfigPath("/usr/etc/")

	flag.String("port", "9680", "The port number to listen on for HTTP requests")
	flag.String("address", "0.0.0.0", "The address to listen on for HTTP requests")
	flag.String("log-level", "info", "The minimum logging level; levels are, in ascending order: debug, info, warn, error")
	flag.String("sap-control-url", "", "The URL of the SAPControl SOAP web service, e.g. http://$HOST:$PORT")

	err := config.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Errorf("Could not bind config to CLI flags: %v", err)
	}
}

func main() {
	initConfig()

	var err error

	client := soap.NewClient(
		config.GetString("sap-control-url"),
		soap.WithBasicAuth(
			config.GetString("sap-control-user"),
			config.GetString("sap-control-password"),
		),
	)
	webService := sapcontrol.NewWebService(client)

	startServiceCollector, err := start_service.NewCollector(webService)

	if err != nil {
		log.Warn(err)
	} else {
		prometheus.MustRegister(startServiceCollector)
		log.Info("Start Service collector registered")
	}

	enqueueServerCollector, err := enqueue_server.NewCollector(webService)

	if err != nil {
		log.Warn(err)
	} else {
		prometheus.MustRegister(enqueueServerCollector)
		log.Info("Enqueue Server collector registered")
	}

	dispatcherCollector, err := dispatcher.NewCollector(webService)

	if err != nil {
		log.Warn(err)
	} else {
		prometheus.MustRegister(dispatcherCollector)
		log.Info("Dispatcher collector registered")
	}

	fullListenAddress := fmt.Sprintf("%s:%s", config.Get("address"), config.Get("port"))

	http.HandleFunc("/", internal.Landing)
	http.Handle("/metrics", promhttp.Handler())

	log.Infof("Serving metrics on %s", fullListenAddress)
	log.Fatal(http.ListenAndServe(fullListenAddress, nil))
}

func initConfig() {
	var err error

	flag.Parse()

	err = config.ReadInConfig()
	if err != nil {
		log.Warn(err)
		log.Info("Default config values will be used")
	} else {
		log.Info("Using config file: ", config.ConfigFileUsed())
	}

	internal.SetLogLevel(config.GetString("log-level"))

	if config.GetString("sap-control-url") == "" {
		log.Fatal("sap-control-url cannot be empty, please use the --sap-control-url flag or set a value in the config")
	}
}
