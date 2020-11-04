package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	// "os/signal"

	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tomi-engel/go-zfs"
	"github.com/tomi-engel/zfs_exporter"
)

var (
	buildVersion = "v0.2.0 (2017.05.16.013)"

	telemetryAddr = flag.String("telemetry.addr", ":9134", "host:port for ZFS exporter")
	telemetryPath = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	poolNames = flag.String("zfs.pools", "", "[optional] comma-separated list of pool names; if none specified, all pools will be scraped")

	// featureZpoolMetricsDisabled = flag.Bool("feature.zpoolMetricsDisabled", false, "set to true if system has the 'zpool get -p' bug")
)

func main() {
	os.Exit(Main())
}

func Main() int {
	flag.Parse()

	var names []string
	if *poolNames == "" {
		var err error
		names, err = zfs.ListZpoolNames()
		if err != nil {
			log.Fatalf("failed to retrieve all ZFS pool names: %v", err)
			return 1
		}
	} else {
		names = strings.Split(*poolNames, ",")
	}

	if len(names) == 0 {
		log.Fatal("no ZFS pools detected, exiting")
		return 1
	}

	// Register our exporter
	prometheus.MustRegister(zfsexporter.New(names))

	// register http handlers
	http.Handle(*telemetryPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *telemetryPath, http.StatusMovedPermanently)
	})

	log.Printf("starting ZFS exporter (%s) on %q for pool(s): %s\n",
		buildVersion, *telemetryAddr, strings.Join(names, ", "))

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("unexpected failure of ZFS exporter HTTP server: %v", err)
		return 1
	}

	return 0
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
