package collector

import (
	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/collector/enqueue_server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func registerOptionalCollectors() {
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
}
