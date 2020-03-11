package enqueue_server

import (
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func NewCollector(webService sapcontrol.WebService) (*enqueueServerCollector, error) {

	c := &enqueueServerCollector{
		collector.NewDefaultCollector("enqueue_server"),
		webService,
	}

	// TO-DO: describe all these metrics
	// https://github.com/SUSE/sap_host_exporter/issues/11
	c.SetDescriptor("owner_now", "TBD", nil)
	c.SetDescriptor("owner_high", "TBD", nil)
	c.SetDescriptor("owner_max", "TBD", nil)
	c.SetDescriptor("owner_state", "TBD", nil)
	c.SetDescriptor("arguments_now", "TBD", nil)
	c.SetDescriptor("arguments_high", "TBD", nil)
	c.SetDescriptor("arguments_max", "TBD", nil)
	c.SetDescriptor("arguments_state", "TBD", nil)
	c.SetDescriptor("locks_now", "TBD", nil)
	c.SetDescriptor("locks_high", "TBD", nil)
	c.SetDescriptor("locks_max", "TBD", nil)
	c.SetDescriptor("locks_state", "TBD", nil)
	c.SetDescriptor("enqueue_requests", "TBD", nil)
	c.SetDescriptor("enqueue_rejects", "TBD", nil)
	c.SetDescriptor("enqueue_errors", "TBD", nil)
	c.SetDescriptor("dequeue_requests", "TBD", nil)
	c.SetDescriptor("dequeue_errors", "TBD", nil)
	c.SetDescriptor("dequeue_all_requests", "TBD", nil)
	c.SetDescriptor("cleanup_requests", "TBD", nil)
	c.SetDescriptor("backup_requests", "TBD", nil)
	c.SetDescriptor("reporting_requests", "TBD", nil)
	c.SetDescriptor("compress_requests", "TBD", nil)
	c.SetDescriptor("verify_requests", "TBD", nil)
	c.SetDescriptor("lock_time", "TBD", nil)
	c.SetDescriptor("lock_wait_time", "TBD", nil)
	c.SetDescriptor("server_time", "TBD", nil)
	c.SetDescriptor("replication_state", "TBD", nil)

	return c, nil
}

type enqueueServerCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *enqueueServerCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting Enqueue Server metrics")

	c.recordEnqStats(ch)
}

func (c *enqueueServerCollector) recordEnqStats(ch chan<- prometheus.Metric) {
	enqStatistic, err := c.webService.EnqGetStatistic()

	if err != nil {
		log.Warnf("SAPControl web service error: %s", err)
		return
	}

	ch <- c.MakeGaugeMetric("owner_now", float64(enqStatistic.OwnerNow))
	ch <- c.MakeGaugeMetric("owner_high", float64(enqStatistic.OwnerHigh))
	ch <- c.MakeGaugeMetric("owner_max", float64(enqStatistic.OwnerMax))

	ownerState, err := sapcontrol.StateColorToFloat(enqStatistic.OwnerState)
	if err != nil {
		log.Warnf("Could not record owner_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("owner_state", ownerState)
	}

	ch <- c.MakeGaugeMetric("arguments_now", float64(enqStatistic.ArgumentsNow))
	ch <- c.MakeGaugeMetric("arguments_high", float64(enqStatistic.ArgumentsHigh))
	ch <- c.MakeGaugeMetric("arguments_max", float64(enqStatistic.ArgumentsMax))

	argumentsState, err := sapcontrol.StateColorToFloat(enqStatistic.ArgumentsState)
	if err != nil {
		log.Warnf("Could not record arguments_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("arguments_state", argumentsState)
	}

	ch <- c.MakeGaugeMetric("locks_now", float64(enqStatistic.LocksNow))
	ch <- c.MakeGaugeMetric("locks_high", float64(enqStatistic.LocksHigh))
	ch <- c.MakeGaugeMetric("locks_max", float64(enqStatistic.LocksMax))

	locksState, err := sapcontrol.StateColorToFloat(enqStatistic.LocksState)
	if err != nil {
		log.Warnf("Could not record locks_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("locks_state", locksState)
	}

	ch <- c.MakeGaugeMetric("enqueue_requests", float64(enqStatistic.EnqueueRequests))
	ch <- c.MakeGaugeMetric("enqueue_rejects", float64(enqStatistic.EnqueueRejects))
	ch <- c.MakeGaugeMetric("enqueue_errors", float64(enqStatistic.EnqueueErrors))

	ch <- c.MakeGaugeMetric("dequeue_requests", float64(enqStatistic.DequeueRequests))
	ch <- c.MakeGaugeMetric("dequeue_errors", float64(enqStatistic.DequeueErrors))
	ch <- c.MakeGaugeMetric("dequeue_all_requests", float64(enqStatistic.DequeueAllRequests))

	ch <- c.MakeGaugeMetric("cleanup_requests", float64(enqStatistic.CleanupRequests))
	ch <- c.MakeGaugeMetric("backup_requests", float64(enqStatistic.BackupRequests))
	ch <- c.MakeGaugeMetric("reporting_requests", float64(enqStatistic.ReportingRequests))
	ch <- c.MakeGaugeMetric("compress_requests", float64(enqStatistic.CompressRequests))
	ch <- c.MakeGaugeMetric("verify_requests", float64(enqStatistic.VerifyRequests))

	ch <- c.MakeGaugeMetric("lock_time", float64(enqStatistic.LockTime))
	ch <- c.MakeGaugeMetric("lock_wait_time", float64(enqStatistic.LockWaitTime))
	ch <- c.MakeGaugeMetric("server_time", float64(enqStatistic.ServerTime))

	replicationState, err := sapcontrol.StateColorToFloat(enqStatistic.ReplicationState)
	if err != nil {
		log.Warnf("Could not record replication_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("replication_state", replicationState)
	}
}
