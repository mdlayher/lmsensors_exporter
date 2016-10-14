package lmsensorsexporter

import (
	"strings"
	"testing"

	"github.com/mdlayher/lmsensors"
)

func TestTemperatureCollector(t *testing.T) {
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
				`lmsensors_temperature_alarm`,
				`lmsensors_temperature_degrees_celsius`,
			},
		},
		{
			name: "one device, one sensor",
			devices: []*lmsensors.Device{{
				Name: "coretemp-00",
				Sensors: []lmsensors.Sensor{
					&lmsensors.TemperatureSensor{
						Name:  "temp1",
						Label: "Core 1",
						Input: 42.0,
					},
				},
			}},
			match: []string{
				`lmsensors_temperature_alarm{details="Core 1",device="coretemp-00",sensor="temp1"} 0`,
				`lmsensors_temperature_degrees_celsius{details="Core 1",device="coretemp-00",sensor="temp1"} 42`,
			},
		},
		{
			name: "two devices, multiple sensors",
			devices: []*lmsensors.Device{
				{
					Name: "coretemp-00",
					Sensors: []lmsensors.Sensor{
						&lmsensors.TemperatureSensor{
							Name:  "temp1",
							Label: "Core 1",
							Input: 42.0,
						},
						&lmsensors.TemperatureSensor{
							Name:  "temp2",
							Label: "Core 2",
							Input: 60.2,
							Alarm: true,
						},
					},
				},
				{
					Name: "it8728-00",
					Sensors: []lmsensors.Sensor{
						&lmsensors.FanSensor{
							Name:  "fan1",
							Input: 1010,
						},
						&lmsensors.TemperatureSensor{
							Name:  "temp1",
							Input: 43.0,
						},
					},
				},
			},
			match: []string{
				`lmsensors_temperature_alarm{details="Core 1",device="coretemp-00",sensor="temp1"} 0`,
				`lmsensors_temperature_alarm{details="Core 2",device="coretemp-00",sensor="temp2"} 1`,
				`lmsensors_temperature_alarm{details="",device="it8728-00",sensor="temp1"} 0`,
				`lmsensors_temperature_degrees_celsius{details="Core 1",device="coretemp-00",sensor="temp1"} 42`,
				`lmsensors_temperature_degrees_celsius{details="Core 2",device="coretemp-00",sensor="temp2"} 60.2`,
				`lmsensors_temperature_degrees_celsius{details="",device="it8728-00",sensor="temp1"} 43`,
			},
			noMatch: []string{
				`lmsensors_fan_rpm`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testCollector(t, NewTemperatureCollector(tt.devices))

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
