// Command lmsensors_exporter provides a Prometheus exporter for lmsensors
// sensor metrics.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mdlayher/lmsensors"
	"github.com/mdlayher/lmsensors_exporter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9165", "address for lmsensors exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(lmsensorsexporter.New(lmsensors.New()))

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting lmsensors exporter on %q", *telemetryAddr)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start lmsensors exporter: %s", err)
	}
}
