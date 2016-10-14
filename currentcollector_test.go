package lmsensorsexporter

import (
	"strings"
	"testing"

	"github.com/mdlayher/lmsensors"
)

func TestCurrentCollector(t *testing.T) {
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
				`lmsensors_current_alarm`,
				`lmsensors_current_amperes`,
			},
		},
		{
			name: "one device, one sensor",
			devices: []*lmsensors.Device{{
				Name: "sfc-00",
				Sensors: []lmsensors.Sensor{
					&lmsensors.CurrentSensor{
						Name:  "curr1",
						Label: "0.9V supply current",
						Input: 7.613,
					},
				},
			}},
			match: []string{
				`lmsensors_current_alarm{details="0.9V supply current",device="sfc-00",sensor="curr1"} 0`,
				`lmsensors_current_amperes{details="0.9V supply current",device="sfc-00",sensor="curr1"} 7.613`,
			},
		},
		{
			name: "two devices, multiple sensors",
			devices: []*lmsensors.Device{
				{
					Name: "sfc-00",
					Sensors: []lmsensors.Sensor{
						&lmsensors.CurrentSensor{
							Name:  "curr1",
							Label: "0.9V supply current",
							Input: 7.613,
						},
						&lmsensors.CurrentSensor{
							Name:  "curr2",
							Label: "",
							Input: 16.123,
							Alarm: true,
						},
					},
				},
				{
					Name: "sfc-01",
					Sensors: []lmsensors.Sensor{
						&lmsensors.FanSensor{
							Name:  "fan1",
							Input: 1010,
						},
						&lmsensors.CurrentSensor{
							Name:  "curr1",
							Input: 3.10,
						},
					},
				},
			},
			match: []string{
				`lmsensors_current_alarm{details="0.9V supply current",device="sfc-00",sensor="curr1"} 0`,
				`lmsensors_current_alarm{details="",device="sfc-00",sensor="curr2"} 1`,
				`lmsensors_current_alarm{details="",device="sfc-01",sensor="curr1"} 0`,
				`lmsensors_current_amperes{details="0.9V supply current",device="sfc-00",sensor="curr1"} 7.613`,
				`lmsensors_current_amperes{details="",device="sfc-00",sensor="curr2"} 16.123`,
				`lmsensors_current_amperes{details="",device="sfc-01",sensor="curr1"} 3.1`,
			},
			noMatch: []string{
				`lmsensors_fan_rpm`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testCollector(t, NewCurrentCollector(tt.devices))

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
