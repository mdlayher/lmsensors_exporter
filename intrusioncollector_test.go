package lmsensorsexporter

import (
	"strings"
	"testing"

	"github.com/mdlayher/lmsensors"
)

func TestIntrusionCollector(t *testing.T) {
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
				`lmsensors_intrusion_alarm`,
			},
		},
		{
			name: "one device, one sensor",
			devices: []*lmsensors.Device{{
				Name: "it8728-00",
				Sensors: []lmsensors.Sensor{
					&lmsensors.IntrusionSensor{
						Name:  "intrusion0",
						Alarm: false,
					},
				},
			}},
			match: []string{
				`lmsensors_intrusion_alarm{device="it8728-00",sensor="intrusion0"} 0`,
			},
		},
		{
			name: "two devices, multiple sensors",
			devices: []*lmsensors.Device{
				{
					Name: "it8728-00",
					Sensors: []lmsensors.Sensor{
						&lmsensors.IntrusionSensor{
							Name:  "intrusion0",
							Alarm: true,
						},
					},
				},
				{
					Name: "it8728-01",
					Sensors: []lmsensors.Sensor{
						&lmsensors.FanSensor{
							Name:  "fan1",
							Input: 998,
						},
						&lmsensors.IntrusionSensor{
							Name:  "intrusion0",
							Alarm: false,
						},
					},
				},
			},
			match: []string{
				`lmsensors_intrusion_alarm{device="it8728-00",sensor="intrusion0"} 1`,
				`lmsensors_intrusion_alarm{device="it8728-01",sensor="intrusion0"} 0`,
			},
			noMatch: []string{
				`lmsensors_fan_rpm`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testCollector(t, NewIntrusionCollector(tt.devices))

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
