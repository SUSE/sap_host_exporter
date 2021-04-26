package enqueue_server

import (
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/lib/sapcontrol"
)

func NewCollector(webService sapcontrol.WebService) (*enqueueServerCollector, error) {

	c := &enqueueServerCollector{
		collector.NewDefaultCollector("enqueue_server"),
		webService,
	}

	c.SetDescriptor("owner_now", "Current number of lock owners in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("owner_high", "Peak number of lock owners that have been stored simultaneously in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("owner_max", "Maximum number of lock owner IDs that can be stored in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("owner_state", "General state of lock owners", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("arguments_now", "Current number of lock arguments in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("arguments_high", "Peak number of lock arguments that have been stored simultaneously in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("arguments_max", "Maximum number of lock arguments that can be stored in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("arguments_state", "General state of lock arguments", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("locks_now", "Current number of elementary locks in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("locks_high", "Peak number of elementary locks that have been stored simultaneously in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("locks_max", "Maximum number of elementary locks that can be stored in the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("locks_state", "General state of elementary locks", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("enqueue_requests", "Lock acquisition requests", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("enqueue_rejects", "Rejected lock requests", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("enqueue_errors", "Lock acquisition errors", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("dequeue_requests", "Lock release requests", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("dequeue_errors", "Lock release errors", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("dequeue_all_requests", "Requests to release of all the locks of an LUW", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("cleanup_requests", "Requests to release of all the locks of an application server", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("backup_requests", "Number of requests forwarded to the update process", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("reporting_requests", "Number of reading operations on the lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("compress_requests", "Internal use", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("verify_requests", "Internal use", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("lock_time", "Total time spent in lock operations", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("lock_wait_time", "Total waiting time of all work processes for accessing lock table", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("server_time", "Total time spent in lock operations by all processes in the enqueue server", []string{"instance_name", "instance_number", "SID", "instance_hostname"})
	c.SetDescriptor("replication_state", "General state of lock server replication", []string{"instance_name", "instance_number", "SID", "instance_hostname"})

	return c, nil
}

type enqueueServerCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *enqueueServerCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting Enqueue Server metrics")

	err := c.recordEnqStats(ch)
	if err != nil {
		log.Warnf("Enqueue Server Collector scrape failed: %s", err)
	}
}

func (c *enqueueServerCollector) recordEnqStats(ch chan<- prometheus.Metric) error {
	enqStatistic, err := c.webService.EnqGetStatistic()
	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	currentSapInstance, err := c.webService.GetCurrentInstance()
	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	labels := []string{
		currentSapInstance.Name,
		strconv.Itoa(int(currentSapInstance.Number)),
		currentSapInstance.SID,
		currentSapInstance.Hostname,
	}

	ch <- c.MakeGaugeMetric("owner_now", float64(enqStatistic.OwnerNow), labels...)
	ch <- c.MakeCounterMetric("owner_high", float64(enqStatistic.OwnerHigh), labels...)
	ch <- c.MakeGaugeMetric("owner_max", float64(enqStatistic.OwnerMax), labels...)

	ownerState, err := sapcontrol.StateColorToFloat(enqStatistic.OwnerState)
	if err != nil {
		log.Warnf("Could not record owner_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("owner_state", ownerState, labels...)
	}

	ch <- c.MakeGaugeMetric("arguments_now", float64(enqStatistic.ArgumentsNow), labels...)
	ch <- c.MakeCounterMetric("arguments_high", float64(enqStatistic.ArgumentsHigh), labels...)
	ch <- c.MakeGaugeMetric("arguments_max", float64(enqStatistic.ArgumentsMax), labels...)

	argumentsState, err := sapcontrol.StateColorToFloat(enqStatistic.ArgumentsState)
	if err != nil {
		log.Warnf("Could not record arguments_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("arguments_state", argumentsState, labels...)
	}

	ch <- c.MakeGaugeMetric("locks_now", float64(enqStatistic.LocksNow), labels...)
	ch <- c.MakeCounterMetric("locks_high", float64(enqStatistic.LocksHigh), labels...)
	ch <- c.MakeGaugeMetric("locks_max", float64(enqStatistic.LocksMax), labels...)

	locksState, err := sapcontrol.StateColorToFloat(enqStatistic.LocksState)
	if err != nil {
		log.Warnf("Could not record locks_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("locks_state", locksState, labels...)
	}

	ch <- c.MakeCounterMetric("enqueue_requests", float64(enqStatistic.EnqueueRequests), labels...)
	ch <- c.MakeCounterMetric("enqueue_rejects", float64(enqStatistic.EnqueueRejects), labels...)
	ch <- c.MakeCounterMetric("enqueue_errors", float64(enqStatistic.EnqueueErrors), labels...)

	ch <- c.MakeCounterMetric("dequeue_requests", float64(enqStatistic.DequeueRequests), labels...)
	ch <- c.MakeCounterMetric("dequeue_errors", float64(enqStatistic.DequeueErrors), labels...)
	ch <- c.MakeCounterMetric("dequeue_all_requests", float64(enqStatistic.DequeueAllRequests), labels...)

	ch <- c.MakeCounterMetric("cleanup_requests", float64(enqStatistic.CleanupRequests), labels...)
	ch <- c.MakeCounterMetric("backup_requests", float64(enqStatistic.BackupRequests), labels...)
	ch <- c.MakeCounterMetric("reporting_requests", float64(enqStatistic.ReportingRequests), labels...)
	ch <- c.MakeCounterMetric("compress_requests", float64(enqStatistic.CompressRequests), labels...)
	ch <- c.MakeCounterMetric("verify_requests", float64(enqStatistic.VerifyRequests), labels...)

	ch <- c.MakeCounterMetric("lock_time", enqStatistic.LockTime, labels...)
	ch <- c.MakeCounterMetric("lock_wait_time", enqStatistic.LockWaitTime, labels...)
	ch <- c.MakeCounterMetric("server_time", enqStatistic.ServerTime, labels...)

	replicationState, err := sapcontrol.StateColorToFloat(enqStatistic.ReplicationState)
	if err != nil {
		log.Warnf("Could not record replication_state metric: %s", err)
	} else {
		ch <- c.MakeGaugeMetric("replication_state", replicationState, labels...)
	}

	return nil
}
