package alert

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SUSE/sap_host_exporter/collector"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
)

func NewCollector(webService sapcontrol.WebService) (*dispatcherCollector, error) {

	c := &dispatcherCollector{
		collector.NewDefaultCollector("alert"),
		webService,
	}

	c.SetDescriptor("ha_check", "High Availability system configuration and status checks", []string{"category", "state", "description", "comment"})
	c.SetDescriptor("ha_failover_active", "Whether or not High Availability Failover is active", nil)

	return c, nil
}

type dispatcherCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *dispatcherCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting Alert metrics")

	var err error

	err = collector.RecordConcurrently([]func(ch chan<- prometheus.Metric) error{
		c.recordHAConfigChecks,
		c.recordHAFailoverConfigChecks,
		c.recordHAFailoverActive,
	}, ch)

	if err != nil {
		log.Warn(err)
	}
}

func (c *dispatcherCollector) recordHAConfigChecks(ch chan<- prometheus.Metric) error {
	response, err := c.webService.HACheckConfig()

	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	err = c.recordHAChecks(response.Checks, ch)
	if err != nil {
		return err
	}

	return nil
}

func (c *dispatcherCollector) recordHAFailoverConfigChecks(ch chan<- prometheus.Metric) error {
	response, err := c.webService.HACheckFailoverConfig()

	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	err = c.recordHAChecks(response.Checks, ch)
	if err != nil {
		return errors.Wrap(err, "Could not record HACheck")
	}

	return nil
}

func (c *dispatcherCollector) recordHAChecks(checks []*sapcontrol.HACheck, ch chan<- prometheus.Metric) error {
	for _, check := range checks {
		err := c.recordHACheck(check, ch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *dispatcherCollector) recordHACheck(check *sapcontrol.HACheck, ch chan<- prometheus.Metric) error {
	stateCode, err := sapcontrol.HaVerificationStateToFloat(check.State)
	state, err := sapcontrol.HaVerificationStateToString(check.State)
	category, err := sapcontrol.HaCheckCategoryToString(check.Category)
	if err != nil {
		return errors.Wrapf(err, "Unable to process SAPControl HACheck data: %v", *check)
	}
	ch <- c.MakeGaugeMetric("ha_check", stateCode, category, state, check.Description, check.Comment)

	return nil
}

func (c *dispatcherCollector) recordHAFailoverActive(ch chan<- prometheus.Metric) error {
	response, err := c.webService.HAGetFailoverConfig()

	if err != nil {
		return errors.Wrap(err, "SAPControl web service error")
	}

	var haActive float64
	if response.HAActive {
		haActive = 1
	}
	ch <- c.MakeGaugeMetric("ha_failover_active", haActive)

	return nil
}
