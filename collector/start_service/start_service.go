package start_service

import (
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func NewCollector(webService sapcontrol.WebService, currentSapInstance sapcontrol.CurrentSapInstance) (*startServiceCollector, error) {

	c := &startServiceCollector{
		collector.NewDefaultCollector("start_service"),
		webService,
		currentSapInstance,
	}

	c.SetDescriptor("processes", "The processes started by the SAP Start Service", []string{"name", "pid", "status", "instance_name", "instance_number", "sid", "instance_hostname"})
	c.SetDescriptor("instances", "All instances of the whole SAP system", []string{"features", "start_priority", "instance_name", "instance_number", "sid", "instance_hostname"})

	return c, nil
}

type startServiceCollector struct {
	collector.DefaultCollector
	webService         sapcontrol.WebService
	currentSapInstance sapcontrol.CurrentSapInstance
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
		ch <- c.MakeGaugeMetric(
			"processes",
			state,
			process.Name,
			strconv.Itoa(int(process.Pid)),
			process.Textstatus,
			c.currentSapInstance.Name,
			strconv.Itoa(int(c.currentSapInstance.Number)),
			c.currentSapInstance.SID,
			c.currentSapInstance.Hostname)
	}

	return nil
}

func (c *startServiceCollector) recordInstances(ch chan<- prometheus.Metric) error {
	instanceList, err := c.webService.GetSystemInstanceList()

	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	for _, instance := range instanceList.Instances {
		// we only record the line relative to the current instance, to avoid duplicated metrics
		// we need to check both instance nr and virtual hostname because with SAP you can never be safe enough
		if instance.InstanceNr != c.currentSapInstance.Number || instance.Hostname != c.currentSapInstance.Hostname {
			continue
		}
		instanceStatus, err := sapcontrol.StateColorToFloat(instance.Dispstatus)
		if err != nil {
			return errors.Wrapf(err, "unable to process SAPControl Instance data: %v", *instance)
		}
		ch <- c.MakeGaugeMetric(
			"instances",
			instanceStatus,
			instance.Features,
			instance.StartPriority,
			c.currentSapInstance.Name,
			strconv.Itoa(int(c.currentSapInstance.Number)),
			c.currentSapInstance.SID,
			c.currentSapInstance.Hostname)
	}

	return nil
}
