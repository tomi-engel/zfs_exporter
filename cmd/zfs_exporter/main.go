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
	telemetryAddr = flag.String("telemetry.addr", ":9134", "host:port for ZFS exporter")
	telemetryPath = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	poolNames = flag.String("zfs.pools", "", "[optional] comma-separated list of pool names; if none specified, all pools will be scraped")
)

func main() {
	flag.Parse()

	var names []string
	if *poolNames == "" {
		var err error
		names, err = detectAllPoolNames()
		if err != nil {
			log.Fatalf("failed to retrieve all ZFS pool names: %v", err)
		}
	} else {
		names = strings.Split(*poolNames, ",")
	}

	if len(names) == 0 {
		log.Fatal("no ZFS pools detected, exiting")
	}

	// Register our exporter
	prometheus.MustRegister(zfsexporter.New(names))

	// register http handlers
	http.Handle(*telemetryPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *telemetryPath, http.StatusMovedPermanently)
	})

	log.Printf("starting ZFS exporter on %q for pool(s): %s\n",
		*telemetryAddr, strings.Join(names, ", "))

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("unexpected failure of ZFS exporter HTTP server: %v", err)
	}
}

func detectAllPoolNames() ([]string, error) {
	pools, err := zfs.ListZpools()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(pools))
	for _, p := range pools {
		names = append(names, p.Name)
	}

	return names, nil
}
