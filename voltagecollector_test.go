package lmsensorsexporter

import (
	"strings"
	"testing"

	"github.com/mdlayher/lmsensors"
)

func TestVoltageCollector(t *testing.T) {
	tests := []struct {
		name    string
		devices []*lmsensors.Device
		match   []string
		noMatch []string
	}{
		{
			name: "no devices",
			match: []string{
				`go_goroutines`,
			},
			noMatch: []string{
				`lmsensors_voltage_alarm`,
				`lmsensors_voltage_volts`,
			},
		},
		{
			name: "one device, one sensor",
			devices: []*lmsensors.Device{{
				Name: "it8728-00",
				Sensors: []lmsensors.Sensor{
					&lmsensors.VoltageSensor{
						Name:  "in1",
						Label: "Vbat",
						Input: 3.29,
					},
				},
			}},
			match: []string{
				`lmsensors_voltage_alarm{details="Vbat",device="it8728-00",sensor="in1"} 0`,
				`lmsensors_voltage_volts{details="Vbat",device="it8728-00",sensor="in1"} 3.29`,
			},
		},
		{
			name: "two devices, multiple sensors",
			devices: []*lmsensors.Device{
				{
					Name: "it8728-00",
					Sensors: []lmsensors.Sensor{
						&lmsensors.VoltageSensor{
							Name:  "in1",
							Label: "Vbat",
							Input: 3.29,
						},
						&lmsensors.VoltageSensor{
							Name:  "in2",
							Label: "3VSB",
							Input: 6.12,
							Alarm: true,
						},
					},
				},
				{
					Name: "it8728-01",
					Sensors: []lmsensors.Sensor{
						&lmsensors.FanSensor{
							Name:  "fan1",
							Input: 1010,
						},
						&lmsensors.VoltageSensor{
							Name:  "in1",
							Input: 3.10,
						},
					},
				},
			},
			match: []string{
				`lmsensors_voltage_alarm{details="Vbat",device="it8728-00",sensor="in1"} 0`,
				`lmsensors_voltage_alarm{details="3VSB",device="it8728-00",sensor="in2"} 1`,
				`lmsensors_voltage_alarm{details="",device="it8728-01",sensor="in1"} 0`,
				`lmsensors_voltage_volts{details="Vbat",device="it8728-00",sensor="in1"} 3.29`,
				`lmsensors_voltage_volts{details="3VSB",device="it8728-00",sensor="in2"} 6.12`,
				`lmsensors_voltage_volts{details="",device="it8728-01",sensor="in1"} 3.1`,
			},
			noMatch: []string{
				`lmsensors_fan_rpm`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testCollector(t, NewVoltageCollector(tt.devices))

			for _, m := range tt.match {
				t.Run("match/"+m, func(t *testing.T) {
					if !strings.Contains(got, m) {
						t.Fatal("output did not contain expected metric")
					}
				})
			}

			for _, m := range tt.noMatch {
				t.Run("nomatch/"+m, func(t *testing.T) {
					if strings.Contains(got, m) {
						t.Fatal("output contains unexpected metric")
					}
				})
			}
		})
	}
}
