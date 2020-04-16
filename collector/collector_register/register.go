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
	enqueueFound := false
	dispatcherFound := false
	processList, err := webService.GetProcessList()
	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	for _, process := range processList.Processes {
		// if we found msg_server on process name we register the Enqueue Server
		if process.Name == "msg_server" && enqueu_found == false {
			enqueueServerCollector, err := enqueue_server.NewCollector(webService)
			if err != nil {
				return errors.Wrap(err, "error registering Enqueue Server collector")
			} else {
				prometheus.MustRegister(enqueueServerCollector)
				log.Info("Enqueue Server optional collector registered")
			}
			enqueu_found = true
		}
		// if we found disp+work on process name we register the dispatcher collector
		if process.Name == "disp+work" && dispatch_found == false {
			dispatcherCollector, err := dispatcher.NewCollector(webService)
			if err != nil {
				return errors.Wrap(err, "error registering Dispatcher collector")
			} else {
				prometheus.MustRegister(dispatcherCollector)
				log.Info("Dispatcher optional collector registered")
			}
			dispatch_found = true
		}

		if enqueu_found == true && dispatch_found == true {
			break
		}
	}

	return nil
}
