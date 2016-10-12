// Package lmsensorsexporter provides the Exporter type used in the
// lmsensors_exporter Prometheus exporter.
package lmsensorsexporter

import (
	"log"

	"github.com/mdlayher/lmsensors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// namespace is the top-level namespace for this lmsensors exporter.
	namespace = "lmsensors"
)

var _ Scanner = &lmsensors.Scanner{}

// A Scanner is a type that can scan for lmsensors devices.  Scanner is
// implemented by *lmsensors.Scanner.
type Scanner interface {
	Scan() ([]*lmsensors.Device, error)
}

// An Exporter is a Prometheus exporter for lmsensors metrics.
// It wraps all lmsensors metrics collectors and provides a single global
// exporter which can serve metrics.
//
// It implements the prometheus.Collector interface in order to register
// with Prometheus.
type Exporter struct {
	s Scanner
}

var _ prometheus.Collector = &Exporter{}

// New creates a new Exporter which collects metrics using the input Scanner.
func New(s Scanner) *Exporter {
	return &Exporter{
		s: s,
	}
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.withCollectors(func(cs []prometheus.Collector) {
		for _, c := range cs {
			c.Describe(ch)
		}
	})
}

// Collect sends the collected metrics from each of the collectors to
// prometheus.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.withCollectors(func(cs []prometheus.Collector) {
		for _, c := range cs {
			c.Collect(ch)
		}
	})
}

// withCollectors starts an lmsensors device scan and creates a set of prometheus
// collectors, invoking the input closure with the collectors to collect metrics.
func (e *Exporter) withCollectors(fn func(cs []prometheus.Collector)) {
	devices, err := e.s.Scan()
	if err != nil {
		log.Printf("[ERROR] error scanning lmsensors devices: %v", err)
		return
	}

	cs := []prometheus.Collector{
		NewFanCollector(devices),
		NewIntrusionCollector(devices),
		NewTemperatureCollector(devices),
		NewVoltageCollector(devices),
	}

	fn(cs)
}

// boolFloat64 converts a boolean into a float64 with a value of 1.0 if true, and
// 0.0 if false.
func boolFloat64(b bool) float64 {
	if !b {
		return 0.0
	}

	return 1.0
}
