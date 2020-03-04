package start_service

import (
	"github.com/hooklift/gowsdl/soap"
	log "github.com/sirupsen/logrus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/prometheus/client_golang/prometheus"
)

func NewCollector(startServiceUrl string) (*soapCollector, error) {
	webService := sapcontrol.NewWebService(soap.NewClient(startServiceUrl))

	c := &soapCollector{
		collector.NewDefaultCollector("start_svc"),
		webService,
	}

	c.SetDescriptor("process", "The SAP processes running on this instance", []string{"name", "status"})

	return c, nil
}

type soapCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *soapCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting SOAP metrics")

	response, err := c.webService.GetProcessList()

	if err != nil {
		log.Warnf("SAPControl service error: %s", err)
		return
	}

	for _, processItem := range response.Process.Item {
		ch <- c.MakeGaugeMetric("process", 1, processItem.Name, string(processItem.Dispstatus))
	}
}
