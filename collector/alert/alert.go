package alert

import (
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

	return c, nil
}

type dispatcherCollector struct {
	collector.DefaultCollector
	webService sapcontrol.WebService
}

func (c *dispatcherCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Collecting Alert metrics")

	c.recordHaChecks(ch)
}

func (c *dispatcherCollector) recordHaChecks(ch chan<- prometheus.Metric) {
	haChecks, err := c.webService.HACheckConfig()

	if err != nil {
		log.Warnf("SAPControl web service error: %s", err)
		return
	}

	for _, check := range haChecks.Check.Item {
		c.recordHaCheck(check, ch)
	}
}

func (c *dispatcherCollector) recordHaCheck(check *sapcontrol.HACheck, ch chan<- prometheus.Metric) {
	stateCode, err := sapcontrol.HaVerificationStateToFloat(check.State)
	if err != nil {
		log.Warnf("Could not record ha_check metric: %s", err)
		return
	}
	state, err := sapcontrol.HaVerificationStateToString(check.State)
	if err != nil {
		log.Warnf("Could not record ha_check metric: %s", err)
		return
	}
	category, err := sapcontrol.HaCheckCategoryToString(check.Category)
	if err != nil {
		log.Warnf("Could not record ha_check metric: %s", err)
		return
	}
	ch <- c.MakeGaugeMetric("ha_check", stateCode, category, state, check.Description, check.Comment)
}
