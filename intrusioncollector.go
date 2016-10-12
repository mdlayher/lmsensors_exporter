package lmsensorsexporter

import (
	"github.com/mdlayher/lmsensors"
	"github.com/prometheus/client_golang/prometheus"
)

// A IntrusionCollector is a Prometheus collector for lmsensors intrusion
// sensor metrics.
type IntrusionCollector struct {
	Alarm *prometheus.Desc

	devices []*lmsensors.Device
}

var _ prometheus.Collector = &IntrusionCollector{}

// NewIntrusionCollector creates a new IntrusionCollector.
func NewIntrusionCollector(devices []*lmsensors.Device) *IntrusionCollector {
	const (
		subsystem = "intrusion"
	)

	var (
		labels = []string{"device", "sensor"}
	)

	return &IntrusionCollector{
		devices: devices,

		Alarm: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "alarm"),
			"Whether or not an intrusion sensor has triggered an alarm (1 - true, 0 - false).",
			labels,
			nil,
		),
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
func (c *IntrusionCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.Alarm,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric created by the IntrusionCollector
// to the provided prometheus Metric channel.
func (c *IntrusionCollector) Collect(ch chan<- prometheus.Metric) {
	for _, d := range c.devices {
		for _, s := range d.Sensors {
			is, ok := s.(*lmsensors.IntrusionSensor)
			if !ok {
				continue
			}

			labels := []string{
				d.Name,
				is.Name,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Alarm,
				prometheus.GaugeValue,
				boolFloat64(is.Alarm),
				labels...,
			)
		}
	}
}
