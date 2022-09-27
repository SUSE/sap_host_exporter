package registry

import (
	"strings"

	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/collector/enqueue_server"
	"github.com/SUSE/sap_host_exporter/lib/sapcontrol"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
		if strings.Contains(process.Name, "msg_server") {
			enqueueFound = true
		}
		if strings.Contains(process.Name, "disp+work") {
			dispatcherFound = true
		}
		if enqueueFound == true && dispatcherFound == true {
			break
		}
	}
	// if we found msg_server on process name we register the Enqueue Server
	if enqueueFound == true {
		enqueueServerCollector, err := enqueue_server.NewCollector(webService)
		if err != nil {
			return errors.Wrap(err, "error registering Enqueue Server collector")
		} else {
			prometheus.MustRegister(enqueueServerCollector)
			log.Info("Enqueue Server optional collector registered")
		}
	}
	// if we found disp+work on process name we register the dispatcher collector
	if dispatcherFound == true {
		dispatcherCollector, err := dispatcher.NewCollector(webService)
		if err != nil {
			return errors.Wrap(err, "error registering Dispatcher collector")
		} else {
			prometheus.MustRegister(dispatcherCollector)
			log.Info("Dispatcher optional collector registered")
		}
	}

	return nil
}
