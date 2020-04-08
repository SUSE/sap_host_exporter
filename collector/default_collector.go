package collector

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const NAMESPACE = "sap"

type DefaultCollector struct {
	subsystem   string
	descriptors map[string]*prometheus.Desc
}

func NewDefaultCollector(subsystem string) DefaultCollector {
	return DefaultCollector{
		subsystem,
		make(map[string]*prometheus.Desc),
	}
}

func (c *DefaultCollector) GetDescriptor(name string) *prometheus.Desc {
	desc, ok := c.descriptors[name]
	if !ok {
		// we hard panic on this because it's most certainly a coding error
		panic(errors.Errorf("undeclared metric '%s'", name))
	}
	return desc
}

// Convenience wrapper around prometheus.NewDesc constructor.
// Stores a metric descriptor with a fully qualified name like `NAMESPACE_subsystem_name`.
// `name` is the last and most relevant part of the metrics Full Qualified Name;
// `help` is the message displayed in the HELP line
// `variableLabels` is a list of labels to declare. Use `nil` to declare no labels.
func (c *DefaultCollector) SetDescriptor(name, help string, variableLabels []string) {
	c.descriptors[name] = prometheus.NewDesc(prometheus.BuildFQName(NAMESPACE, c.subsystem, name), help, variableLabels, nil)
}

func (c *DefaultCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descriptor := range c.descriptors {
		ch <- descriptor
	}
}

func (c *DefaultCollector) MakeGaugeMetric(name string, value float64, labelValues ...string) prometheus.Metric {
	return c.makeMetric(name, value, prometheus.GaugeValue, labelValues...)
}

func (c *DefaultCollector) MakeCounterMetric(name string, value float64, labelValues ...string) prometheus.Metric {
	return c.makeMetric(name, value, prometheus.CounterValue, labelValues...)
}

func (c *DefaultCollector) makeMetric(name string, value float64, valueType prometheus.ValueType, labelValues ...string) prometheus.Metric {
	desc := c.GetDescriptor(name)
	return prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)
}

// Run multiple metric recording functions concurrently
func RecordConcurrently(recorders []func(ch chan<- prometheus.Metric) error, ch chan<- prometheus.Metric) []error {
	results := make(chan error, len(recorders))
	var errs []error
	var wg sync.WaitGroup

	// For each recorder we start a goroutine which will send its result in a channel.
	// A Waitgroup is used to later wait for all of them.
	for _, recorder := range recorders {
		wg.Add(1)
		go func(recorder func(ch chan<- prometheus.Metric) error, wg *sync.WaitGroup) {
			defer wg.Done()
			results <- recorder(ch)
		}(recorder, &wg)
	}

	// As soon as all the goroutines in the Waitgroup are done, close the channel where the errors are sent
	go func() {
		wg.Wait()
		close(results)
	}()

	// Scroll the results channel and store potential errors in an array. This will block until the channel is closed.
	for err := range results {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
