package start_service

import (
	"strconv"

	"github.com/hooklift/gowsdl/soap"
	log "github.com/sirupsen/logrus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/prometheus/client_golang/prometheus"
)

func NewCollector(startServiceUrl string) (*soapCollector, error) {
	webService := sapcontrol.NewWebService(soap.NewClient(startServiceUrl))

	c := &soapCollector{
		collector.NewDefaultCollector("start_service"),
		webService,
	}

	c.SetDescriptor("processes", "The processes started by the SAP Start Service", []string{"name", "pid", "textstatus", "dispstatus"})

	return c, nil
}

type soapCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *soapCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting SAP Start Service metrics")

	response, err := c.webService.GetProcessList()

	if err != nil {
		log.Warnf("SAPControl web service error: %s", err)
		return
	}

	for _, process := range response.Process.Item {
		dispStatus, _ := sapcontrol.StateColorToString(process.Dispstatus)
		ch <- c.MakeGaugeMetric("processes", 1, process.Name, strconv.Itoa(int(process.Pid)), process.Textstatus, dispStatus)
	}
}
