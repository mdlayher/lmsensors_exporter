package lmsensorsexporter

import (
	"github.com/mdlayher/lmsensors"
	"github.com/prometheus/client_golang/prometheus"
)

// A CurrentCollector is a Prometheus collector for lmsensors current
// sensor metrics.
type CurrentCollector struct {
	Amperes *prometheus.Desc
	Alarm   *prometheus.Desc

	devices []*lmsensors.Device
}

var _ prometheus.Collector = &CurrentCollector{}

// NewCurrentCollector creates a new CurrentCollector.
func NewCurrentCollector(devices []*lmsensors.Device) *CurrentCollector {
	const (
		subsystem = "current"
	)

	var (
		labels = []string{"device", "sensor", "details"}
	)

	return &CurrentCollector{
		devices: devices,

		Amperes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "amperes"),
			"Current current detected by sensor in Amperes.",
			labels,
			nil,
		),

		Alarm: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "alarm"),
			"Whether or not a current sensor has triggered an alarm (1 - true, 0 - false).",
			labels,
			nil,
		),
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
func (c *CurrentCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.Amperes,
		c.Alarm,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric created by the CurrentCollector
// to the provided prometheus Metric channel.
func (c *CurrentCollector) Collect(ch chan<- prometheus.Metric) {
	for _, d := range c.devices {
		for _, s := range d.Sensors {
			cs, ok := s.(*lmsensors.CurrentSensor)
			if !ok {
				continue
			}

			labels := []string{
				d.Name,
				cs.Name,
				cs.Label,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Amperes,
				prometheus.GaugeValue,
				cs.Input,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Alarm,
				prometheus.GaugeValue,
				boolFloat64(cs.Alarm),
				labels...,
			)
		}
	}
}
