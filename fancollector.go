package lmsensorsexporter

import (
	"github.com/mdlayher/lmsensors"
	"github.com/prometheus/client_golang/prometheus"
)

// A FanCollector is a Prometheus collector for lmsensors fan
// sensor metrics.
type FanCollector struct {
	RPM   *prometheus.Desc
	Alarm *prometheus.Desc

	devices []*lmsensors.Device
}

var _ prometheus.Collector = &FanCollector{}

// NewFanCollector creates a new FanCollector.
func NewFanCollector(devices []*lmsensors.Device) *FanCollector {
	const (
		subsystem = "fan"
	)

	var (
		labels = []string{"device", "sensor"}
	)

	return &FanCollector{
		devices: devices,

		RPM: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "rpm"),
			"Current fan speed detected by sensor in rotations per minute.",
			labels,
			nil,
		),

		Alarm: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "alarm"),
			"Whether or not a fan sensor has triggered an alarm (1 - true, 0 - false).",
			labels,
			nil,
		),
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
func (c *FanCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.RPM,
		c.Alarm,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric created by the FanCollector
// to the provided prometheus Metric channel.
func (c *FanCollector) Collect(ch chan<- prometheus.Metric) {
	for _, d := range c.devices {
		for _, s := range d.Sensors {
			fs, ok := s.(*lmsensors.FanSensor)
			if !ok {
				continue
			}

			labels := []string{
				d.Name,
				fs.Name,
			}

			ch <- prometheus.MustNewConstMetric(
				c.RPM,
				prometheus.GaugeValue,
				float64(fs.Input),
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Alarm,
				prometheus.GaugeValue,
				boolFloat64(fs.Alarm),
				labels...,
			)
		}
	}
}
