package dispatcher

import (
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func NewCollector(webService sapcontrol.WebService) (*dispatcherCollector, error) {

	c := &dispatcherCollector{
		collector.NewDefaultCollector("dispatcher"),
		webService,
	}

	c.SetDescriptor("queue_now", "Work process current queue length", []string{"type"})
	c.SetDescriptor("queue_high", "Work process peak queue length", []string{"type"})
	c.SetDescriptor("queue_max", "Work process maximum queue length", []string{"type"})
	c.SetDescriptor("queue_writes", "Work process queue writes", []string{"type"})
	c.SetDescriptor("queue_reads", "Work process queue reads", []string{"type"})

	return c, nil
}

type dispatcherCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *dispatcherCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting Dispatcher metrics")

	c.recordWorkProcessQueueStats(ch)
}

func (c *dispatcherCollector) recordWorkProcessQueueStats(ch chan<- prometheus.Metric) {
	queueStatistic, err := c.webService.GetQueueStatistic()

	if err != nil {
		log.Warnf("SAPControl web service error: %s", err)
		return
	}

	// for each work queue, we record a different line for each stat of that queue, with the type as a common label
	for _, queue := range queueStatistic.Queues {
		ch <- c.MakeGaugeMetric("queue_now", float64(queue.Now), queue.Type)
		ch <- c.MakeCounterMetric("queue_high", float64(queue.High), queue.Type)
		ch <- c.MakeGaugeMetric("queue_max", float64(queue.Max), queue.Type)
		ch <- c.MakeCounterMetric("queue_writes", float64(queue.Writes), queue.Type)
		ch <- c.MakeCounterMetric("queue_reads", float64(queue.Reads), queue.Type)
	}
}
