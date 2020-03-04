package main

import (
	"fmt"
	"net/http"

	_ "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	config "github.com/spf13/viper"

	"github.com/SUSE/sap_host_exporter/internal"
	"github.com/SUSE/sap_host_exporter/sapcontrol"
	"github.com/davecgh/go-spew/spew"

	"github.com/hooklift/gowsdl/soap"
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

	err := config.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Errorf("Could not bind config to CLI flags: %v", err)
	}
}

func main() {
	var err error

	flag.Parse()

	err = config.ReadInConfig()
	if err != nil {
		log.Warn(err)
		log.Info("Default config values will be used")
	} else {
		log.Info("Using config file: ", config.ConfigFileUsed())
	}

	// call prometheus.MustRegister() here after instantiating collectors

	internal.SetLogLevel(config.GetString("log-level"))

	fullListenAddress := fmt.Sprintf("%s:%s", config.Get("address"), config.Get("port"))

	http.HandleFunc("/", internal.Landing)
	http.Handle("/metrics", promhttp.Handler())

	client := soap.NewClient("http://10.162.32.183:50013")
	service := sapcontrol.NewSAPControlPortType(client)

	reply, err := service.GetProcessList(&sapcontrol.GetProcessList{})
	if err != nil {
		log.Fatalf("could't get trade prices: %v", err)
	}
	spew.Dump(reply)

	log.Infof("Serving metrics on %s", fullListenAddress)
	log.Fatal(http.ListenAndServe(fullListenAddress, nil))
}
