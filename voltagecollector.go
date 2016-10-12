package lmsensorsexporter

import (
	"github.com/mdlayher/lmsensors"
	"github.com/prometheus/client_golang/prometheus"
)

// A VoltageCollector is a Prometheus collector for lmsensors voltage
// sensor metrics.
type VoltageCollector struct {
	Volts *prometheus.Desc
	Alarm *prometheus.Desc

	devices []*lmsensors.Device
}

var _ prometheus.Collector = &VoltageCollector{}

// NewVoltageCollector creates a new VoltageCollector.
func NewVoltageCollector(devices []*lmsensors.Device) *VoltageCollector {
	const (
		subsystem = "voltage"
	)

	var (
		labels = []string{"device", "sensor", "details"}
	)

	return &VoltageCollector{
		devices: devices,

		Volts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "volts"),
			"Current voltage detected by sensor in volts.",
			labels,
			nil,
		),

		Alarm: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "alarm"),
			"Whether or not a voltage sensor has triggered an alarm (1 - true, 0 - false).",
			labels,
			nil,
		),
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
func (c *VoltageCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.Volts,
		c.Alarm,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric created by the VoltageCollector
// to the provided prometheus Metric channel.
func (c *VoltageCollector) Collect(ch chan<- prometheus.Metric) {
	for _, d := range c.devices {
		for _, s := range d.Sensors {
			vs, ok := s.(*lmsensors.VoltageSensor)
			if !ok {
				continue
			}

			labels := []string{
				d.Name,
				vs.Name,
				vs.Label,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Volts,
				prometheus.GaugeValue,
				vs.Current,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Alarm,
				prometheus.GaugeValue,
				boolFloat64(vs.Alarm),
				labels...,
			)
		}
	}
}
