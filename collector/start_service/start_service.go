package start_service

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/soap"
)

func NewCollector() (*soapCollector, error) {

	c := &soapCollector{
		collector.NewDefaultCollector("soap"),
	}

	c.SetDescriptor("process", "get status of the Process running on the host", []string{"status"})

	return c, nil
}

type soapCollector struct {
	collector.DefaultCollector
}

func (c *soapCollector) Collect(ch chan<- prometheus.Metric) {
	log.Infoln("Collecting SOAP metrics")
	soap.Hello()
}
