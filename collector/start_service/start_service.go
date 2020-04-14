package start_service

import (
	"strconv"

	"github.com/pkg/errors"
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

	c.SetDescriptor("processes", "The processes started by the SAP Start Service", []string{"name", "pid", "status"})
	c.SetDescriptor("instances", "All instances of the whole SAP system", []string{"hostname", "instance_number", "start_priority", "features"})

	return c, nil
}

type startServiceCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *startServiceCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting SAP Start Service metrics")

	errs := collector.RecordConcurrently([]func(ch chan<- prometheus.Metric) error{
		c.recordProcesses,
		c.recordInstances,
	}, ch)

	for _, err := range errs {
		log.Warnf("Start Service Collector scrape failed: %s", err)
	}
}

func (c *startServiceCollector) recordProcesses(ch chan<- prometheus.Metric) error {
	processList, err := c.webService.GetProcessList()

	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	for _, process := range processList.Processes {
		state, err := sapcontrol.StateColorToFloat(process.Dispstatus)
		if err != nil {
			return errors.Wrapf(err, "unable to process SAPControl OSProcess data: %v", *process)
		}
		ch <- c.MakeGaugeMetric("processes", state, process.Name, strconv.Itoa(int(process.Pid)), process.Textstatus)
	}

	return nil
}

func (c *startServiceCollector) recordInstances(ch chan<- prometheus.Metric) error {
	instanceList, err := c.webService.GetSystemInstanceList()

	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	for _, instance := range instanceList.Instances {
		state, err := sapcontrol.StateColorToFloat(instance.Dispstatus)
		if err != nil {
			return errors.Wrapf(err, "unable to process SAPControl Instance data: %v", *instance)
		}
		ch <- c.MakeGaugeMetric("instances", state, instance.Hostname, strconv.Itoa(int(instance.InstanceNr)), instance.StartPriority, instance.Features)
	}

	return nil
}
