package collector_register

import (
	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/collector/enqueue_server"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// RegisterOptionalCollectors register depending on the system where the exporter run the additional collectors
func RegisterOptionalCollectors(webService sapcontrol.WebService) {
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
