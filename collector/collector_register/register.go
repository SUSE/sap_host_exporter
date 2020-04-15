package collector_register

import (
	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/collector/enqueue_server"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// RegisterOptionalCollectors register depending on the system where the exporter run the additional collectors
func RegisterOptionalCollectors(webService sapcontrol.WebService) error {
	enqueu_found := false
	dispatch_found := false
	processList, err := webService.GetProcessList()
	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	for _, process := range processList.Processes {
		// if we found msg_server on process name we register the Enqueue Server
		if process.Name == "msg_server" {
			enqueueServerCollector, err := enqueue_server.NewCollector(webService)
			if err != nil {
				return errors.Wrap(err, "Error by registering enqueueServer collector")
			} else {
				prometheus.MustRegister(enqueueServerCollector)
				log.Info("Enqueue Server collector registered")
			}
			enqueu_found = true
		}
		// if we found disp+work on process name we register the dispatcher collector
		if process.Name == "disp+work" {
			dispatcherCollector, err := dispatcher.NewCollector(webService)
			if err != nil {
				return errors.Wrap(err, "Error by registering dispatcher collector")
			} else {
				prometheus.MustRegister(dispatcherCollector)
				log.Info("Dispatcher collector registered")
			}
			dispatch_found = true
		}

		if enqueu_found == true && dispatch_found == true {
			break
		}
	}

	return nil
}
