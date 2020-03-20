package main

import (
	"fmt"
	"net/http"

	"github.com/hooklift/gowsdl/soap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/SUSE/sap_host_exporter/collector/alert"
	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/collector/enqueue_server"
	"github.com/SUSE/sap_host_exporter/collector/start_service"
	"github.com/SUSE/sap_host_exporter/internal"
	"github.com/SUSE/sap_host_exporter/internal/config"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func init() {
	flag.String("port", "9680", "The port number to listen on for HTTP requests")
	flag.String("address", "0.0.0.0", "The address to listen on for HTTP requests")
	flag.String("log-level", "info", "The minimum logging level; levels are, in ascending order: debug, info, warn, error")
	flag.String("sap-control-url", "localhost:50013", "The URL of the SAPControl SOAP web service, e.g. $HOST:$PORT")
	flag.StringP("config", "c", "", "The path to a custom configuration file. NOTE: it must be in yaml format.")
}

func main() {
	var err error

	flag.Parse()

	config, err := config.New()
	if err != nil {
		log.Fatalf("Could not initialize config: %s", err)
	}

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

	alertCollector, err := alert.NewCollector(webService)
	if err != nil {
		log.Warn(err)
	} else {
		prometheus.MustRegister(alertCollector)
		log.Info("Alert collector registered")
	}

	// if we're not in debug log level, we unregister the Go runtime metrics collector that gets registered by default
	if !log.IsLevelEnabled(log.DebugLevel) {
		prometheus.Unregister(prometheus.NewGoCollector())
	}

	fullListenAddress := fmt.Sprintf("%s:%s", config.Get("address"), config.Get("port"))

	http.HandleFunc("/", internal.Landing)
	http.Handle("/metrics", promhttp.Handler())

	log.Infof("Serving metrics on %s", fullListenAddress)
	log.Fatal(http.ListenAndServe(fullListenAddress, nil))
}
