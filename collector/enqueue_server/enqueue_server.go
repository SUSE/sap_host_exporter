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

	c.SetDescriptor("owner_now", "Current number of lock owners in the lock table", nil)
	c.SetDescriptor("owner_high", "Peak number of lock owners that have been stored simultaneously in the lock table", nil)
	c.SetDescriptor("owner_max", "Maximum number of lock owner IDs that can be stored in the lock table", nil)
	c.SetDescriptor("owner_state", "General state of lock owners", nil)
	c.SetDescriptor("arguments_now", "Current number of lock arguments in the lock table", nil)
	c.SetDescriptor("arguments_high", "Peak number of lock arguments that have been stored simultaneously in the lock table", nil)
	c.SetDescriptor("arguments_max", "Maximum number of lock arguments that can be stored in the lock table", nil)
	c.SetDescriptor("arguments_state", "General state of lock arguments", nil)
	c.SetDescriptor("locks_now", "Current number of elementary locks in the lock table", nil)
	c.SetDescriptor("locks_high", "Peak number of elementary locks that have been stored simultaneously in the lock table", nil)
	c.SetDescriptor("locks_max", "Maximum number of elementary locks that can be stored in the lock table", nil)
	c.SetDescriptor("locks_state", "General state of elementary locks", nil)
	c.SetDescriptor("enqueue_requests", "Lock acquisition requests", nil)
	c.SetDescriptor("enqueue_rejects", "Rejected lock requests", nil)
	c.SetDescriptor("enqueue_errors", "Lock acquisition errors", nil)
	c.SetDescriptor("dequeue_requests", "Lock release requests", nil)
	c.SetDescriptor("dequeue_errors", "Lock release errors", nil)
	c.SetDescriptor("dequeue_all_requests", "Requests to release of all the locks of an LUW", nil)
	c.SetDescriptor("cleanup_requests", "Requests to release of all the locks of an application server", nil)
	c.SetDescriptor("backup_requests", "Number of requests forwarded to the update process", nil)
	c.SetDescriptor("reporting_requests", "Number of reading operations on the lock table", nil)
	c.SetDescriptor("compress_requests", "Internal use", nil)
	c.SetDescriptor("verify_requests", "Internal use", nil)
	c.SetDescriptor("lock_time", "Total time spent in lock operations", nil)
	c.SetDescriptor("lock_wait_time", "Total waiting time of all work processes for accessing lock table", nil)
	c.SetDescriptor("server_time", "Total time spent in lock operations by all processes in the enqueue server", nil)
	c.SetDescriptor("replication_state", "General state of lock server replication", nil)

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
	ch <- c.MakeCounterMetric("owner_high", float64(enqStatistic.OwnerHigh))
	ch <- c.MakeGaugeMetric("owner_max", float64(enqStatistic.OwnerMax))

	ownerState, err := sapcontrol.StateColorToFloat(enqStatistic.OwnerState)
	if err != nil {
		log.Warnf("Could not record owner_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("owner_state", ownerState)
	}

	ch <- c.MakeGaugeMetric("arguments_now", float64(enqStatistic.ArgumentsNow))
	ch <- c.MakeCounterMetric("arguments_high", float64(enqStatistic.ArgumentsHigh))
	ch <- c.MakeGaugeMetric("arguments_max", float64(enqStatistic.ArgumentsMax))

	argumentsState, err := sapcontrol.StateColorToFloat(enqStatistic.ArgumentsState)
	if err != nil {
		log.Warnf("Could not record arguments_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("arguments_state", argumentsState)
	}

	ch <- c.MakeGaugeMetric("locks_now", float64(enqStatistic.LocksNow))
	ch <- c.MakeCounterMetric("locks_high", float64(enqStatistic.LocksHigh))
	ch <- c.MakeGaugeMetric("locks_max", float64(enqStatistic.LocksMax))

	locksState, err := sapcontrol.StateColorToFloat(enqStatistic.LocksState)
	if err != nil {
		log.Warnf("Could not record locks_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("locks_state", locksState)
	}

	ch <- c.MakeCounterMetric("enqueue_requests", float64(enqStatistic.EnqueueRequests))
	ch <- c.MakeCounterMetric("enqueue_rejects", float64(enqStatistic.EnqueueRejects))
	ch <- c.MakeCounterMetric("enqueue_errors", float64(enqStatistic.EnqueueErrors))

	ch <- c.MakeCounterMetric("dequeue_requests", float64(enqStatistic.DequeueRequests))
	ch <- c.MakeCounterMetric("dequeue_errors", float64(enqStatistic.DequeueErrors))
	ch <- c.MakeCounterMetric("dequeue_all_requests", float64(enqStatistic.DequeueAllRequests))

	ch <- c.MakeCounterMetric("cleanup_requests", float64(enqStatistic.CleanupRequests))
	ch <- c.MakeCounterMetric("backup_requests", float64(enqStatistic.BackupRequests))
	ch <- c.MakeCounterMetric("reporting_requests", float64(enqStatistic.ReportingRequests))
	ch <- c.MakeCounterMetric("compress_requests", float64(enqStatistic.CompressRequests))
	ch <- c.MakeCounterMetric("verify_requests", float64(enqStatistic.VerifyRequests))

	ch <- c.MakeCounterMetric("lock_time", enqStatistic.LockTime)
	ch <- c.MakeCounterMetric("lock_wait_time", enqStatistic.LockWaitTime)
	ch <- c.MakeCounterMetric("server_time", enqStatistic.ServerTime)

	replicationState, err := sapcontrol.StateColorToFloat(enqStatistic.ReplicationState)
	if err != nil {
		log.Warnf("Could not record replication_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("replication_state", replicationState)
	}
}
