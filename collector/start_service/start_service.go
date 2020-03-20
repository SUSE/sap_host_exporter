package start_service

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func NewCollector(webService sapcontrol.WebService) (*startServiceCollector, error) {

	c := &startServiceCollector{
		collector.NewDefaultCollector("start_service"),
		webService,
	}

	c.SetDescriptor("processes", "The processes started by the SAP Start Service", []string{"name", "pid", "textstatus", "dispstatus"})

	return c, nil
}

type startServiceCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *startServiceCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting SAP Start Service metrics")

	c.recordProcesses(ch)
}

func (c *startServiceCollector) recordProcesses(ch chan<- prometheus.Metric) {
	processList, err := c.webService.GetProcessList()

	if err != nil {
		log.Warnf("SAPControl web service error: %s", err)
		return
	}

	for _, process := range processList.Processes {
		dispStatus, _ := sapcontrol.StateColorToString(process.Dispstatus)
		ch <- c.MakeGaugeMetric("processes", 1, process.Name, strconv.Itoa(int(process.Pid)), process.Textstatus, dispStatus)
	}
}
