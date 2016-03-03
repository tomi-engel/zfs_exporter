package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/eliothedeman/zfs_exporter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9134", "host:port for zfs exporter")
	telemetryPath = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")
)

func main() {

	// Register our exporter
	prometheus.MustRegister(zfsexporter.New())

	// register http handlers
	http.Handle(*telemetryPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *telemetryPath, http.StatusMovedPermanently)
	})

	log.Printf("Starting ZFS exporter on %q", *telemetryAddr)

	// run
	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("Unexpected failure of ZFS exporter: %s", err)
	}

}
