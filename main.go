package main

import (
	"fmt"
	"net/http"

	"github.com/SUSE/sap_host_exporter/collector/registry"
	"github.com/SUSE/sap_host_exporter/collector/start_service"
	"github.com/SUSE/sap_host_exporter/internal"
	"github.com/SUSE/sap_host_exporter/internal/config"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func init() {
	flag.String("port", "9680", "The port number to listen on for HTTP requests")
	flag.String("address", "0.0.0.0", "The address to listen on for HTTP requests")
	flag.String("log-level", "info", "The minimum logging level; levels are, in ascending order: debug, info, warn, error")
	flag.String("sap-control-url", "localhost:50013", "The URL of the SAPControl SOAP web service, e.g. $HOST:$PORT")
	flag.String("sap-control-uds", "", "The path to the SAPControl Unix Domain Socket. If set, this will be used instead of the URL.")
	flag.StringP("config", "c", "", "The path to a custom configuration file. NOTE: it must be in yaml format.")
}

func main() {
	var err error

	flag.Parse()

	config, err := config.New()
	if err != nil {
		log.Fatalf("Could not initialize config: %s", err)
	}

	client := sapcontrol.NewSoapClient(config)
	webService := sapcontrol.NewWebService(client)
	currentSapInstance, err := webService.GetCurrentInstance()
	if err != nil {
		log.Fatal(errors.Wrap(err, "SAPControl web service error"))
	}

	log.Infof("Monitoring SAP Instance %s", currentSapInstance)

	startServiceCollector, err := start_service.NewCollector(webService)
	if err != nil {
		log.Warn(err)
	} else {
		prometheus.MustRegister(startServiceCollector)
		log.Info("Start Service collector registered")
	}

	err = registry.RegisterOptionalCollectors(webService)
	if err != nil {
		log.Fatal(err)
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
