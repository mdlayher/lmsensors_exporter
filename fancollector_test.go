package lmsensorsexporter

import (
	"strings"
	"testing"

	"github.com/mdlayher/lmsensors"
)

func TestFanCollector(t *testing.T) {
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
				`lmsensors_fan_alarm`,
				`lmsensors_fan_rpm`,
			},
		},
		{
			name: "one device, one sensor",
			devices: []*lmsensors.Device{{
				Name: "it8728-00",
				Sensors: []lmsensors.Sensor{
					&lmsensors.FanSensor{
						Name:    "fan1",
						Current: 1000,
					},
				},
			}},
			match: []string{
				`lmsensors_fan_alarm{device="it8728-00",sensor="fan1"} 0`,
				`lmsensors_fan_rpm{device="it8728-00",sensor="fan1"} 1000`,
			},
		},
		{
			name: "two devices, multiple sensors",
			devices: []*lmsensors.Device{
				{
					Name: "it8728-00",
					Sensors: []lmsensors.Sensor{
						&lmsensors.FanSensor{
							Name:    "fan1",
							Current: 1010,
						},
						&lmsensors.FanSensor{
							Name:    "fan2",
							Current: 0,
							Alarm:   true,
						},
					},
				},
				{
					Name: "it8728-01",
					Sensors: []lmsensors.Sensor{
						&lmsensors.FanSensor{
							Name:    "fan1",
							Current: 998,
						},
						&lmsensors.IntrusionSensor{
							Name:  "intrusion0",
							Alarm: true,
						},
					},
				},
			},
			match: []string{
				`lmsensors_fan_alarm{device="it8728-00",sensor="fan1"} 0`,
				`lmsensors_fan_alarm{device="it8728-00",sensor="fan2"} 1`,
				`lmsensors_fan_alarm{device="it8728-01",sensor="fan1"} 0`,
				`lmsensors_fan_rpm{device="it8728-00",sensor="fan1"} 1010`,
				`lmsensors_fan_rpm{device="it8728-00",sensor="fan2"} 0`,
				`lmsensors_fan_rpm{device="it8728-01",sensor="fan1"} 998`,
			},
			noMatch: []string{
				`lmsensors_intrusion_alarm`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testCollector(t, NewFanCollector(tt.devices))

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
