package collector

import (
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetricFactory(t *testing.T) {
	SUT := NewDefaultCollector("test")
	SUT.SetDescriptor("test_metric", "", nil)

	metric := SUT.MakeGaugeMetric("test_metric", 1)

	assert.Equal(t, SUT.GetDescriptor("test_metric"), metric.Desc())
}

func TestRecordConcurrently(t *testing.T) {
	metrics := make(chan prometheus.Metric, 2)
	metric1 := prometheus.NewGauge(prometheus.GaugeOpts{})
	metric2 := prometheus.NewGauge(prometheus.GaugeOpts{})
	recorder1 := func(ch chan<- prometheus.Metric) error {
		// we make metric1 take longer so that we can assert that metric2 will come first
		time.Sleep(time.Millisecond * 50)
		ch <- metric1
		return nil
	}
	recorder2 := func(ch chan<- prometheus.Metric) error {
		ch <- metric2
		return nil
	}

	err := RecordConcurrently([]func(ch chan<- prometheus.Metric) error{recorder1, recorder2}, metrics)
	assert.NoError(t, err)
	assert.Equal(t, metric2, <-metrics)
	assert.Equal(t, metric1, <-metrics)
}

func TestRecordConcurrentlyErrors(t *testing.T) {
	metrics := make(chan prometheus.Metric)
	expectedError := errors.New("")
	recorder1 := func(ch chan<- prometheus.Metric) error {
		return expectedError
	}
	recorder2 := func(ch chan<- prometheus.Metric) error {
		// we make metric1 take longer so that we can assert that metric2 will come first
		time.Sleep(time.Millisecond * 50)
		return nil
	}

	err := RecordConcurrently([]func(ch chan<- prometheus.Metric) error{recorder1, recorder2}, metrics)
	assert.Equal(t, expectedError, err)
}
