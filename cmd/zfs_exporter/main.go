package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/eliothedeman/go-zfs"
	"github.com/eliothedeman/zfs_exporter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9134", "host:port for zfs exporter")
	telemetryPath = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")
	poolNames     = flag.String("pool.names", "", "Comma seperated list of pool names (Defaults to '' which fetches all pool names)")
)

func detectAllPoolNames() []string {
	pools, err := zfs.ListZpools()
	if err != nil {
		log.Printf("Unable to list zpools for this host: %s", err)
		return nil
	}
	names := make([]string, len(pools))
	for i, p := range pools {
		names[i] = p.Name
	}

	return names
}

func main() {
	flag.Parse()

	var names []string
	if *poolNames == "" {

		names = detectAllPoolNames()

		// none found, exit
		if names == nil {
			return
		}

		log.Printf("No pool names provided, using all pools\n")
	} else {
		names = strings.Split(*poolNames, ",")
	}

	log.Printf("Monitoring pools: %v", names)

	// Register our exporter
	prometheus.MustRegister(zfsexporter.New(names))

	// register http handlers
	http.Handle(*telemetryPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *telemetryPath, http.StatusMovedPermanently)
	})

	log.Printf("Starting ZFS exporter on %q\n", *telemetryAddr)

	// run
	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("Unexpected failure of ZFS exporter: %s\n", err)
	}

}
