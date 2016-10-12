package lmsensorsexporter

import (
	"github.com/mdlayher/lmsensors"
	"github.com/prometheus/client_golang/prometheus"
)

// A TemperatureCollector is a Prometheus collector for lmsensors temperature
// sensor metrics.
type TemperatureCollector struct {
	DegreesCelsius *prometheus.Desc
	Alarm          *prometheus.Desc

	devices []*lmsensors.Device
}

var _ prometheus.Collector = &TemperatureCollector{}

// NewTemperatureCollector creates a new TemperatureCollector.
func NewTemperatureCollector(devices []*lmsensors.Device) *TemperatureCollector {
	const (
		subsystem = "temperature"
	)

	var (
		labels = []string{"device", "sensor", "details"}
	)

	return &TemperatureCollector{
		devices: devices,

		DegreesCelsius: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "degrees_celsius"),
			"Current temperature detected by sensor in degrees Celsius.",
			labels,
			nil,
		),

		Alarm: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "alarm"),
			"Whether or not a temperature sensor has triggered an alarm (1 - true, 0 - false).",
			labels,
			nil,
		),
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
func (c *TemperatureCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.DegreesCelsius,
		c.Alarm,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric created by the TemperatureCollector
// to the provided prometheus Metric channel.
func (c *TemperatureCollector) Collect(ch chan<- prometheus.Metric) {
	for _, d := range c.devices {
		for _, s := range d.Sensors {
			ts, ok := s.(*lmsensors.TemperatureSensor)
			if !ok {
				continue
			}

			labels := []string{
				d.Name,
				ts.Name,
				ts.Label,
			}

			ch <- prometheus.MustNewConstMetric(
				c.DegreesCelsius,
				prometheus.GaugeValue,
				ts.Current,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Alarm,
				prometheus.GaugeValue,
				boolFloat64(ts.Alarm),
				labels...,
			)
		}
	}
}
