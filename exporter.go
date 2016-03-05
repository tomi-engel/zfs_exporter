package zfsexporter

import (
	"log"
	"sync"

	"github.com/eliothedeman/go-zfs"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "zfs"
)

// An Exporter is a Prometheus exporter for the ZFS filesystem.
type Exporter struct {
	mu        sync.Mutex
	poolNames []string
}

// make sure the exporter satisfies the prometheus collector interface
var _ prometheus.Collector = &Exporter{}

// New creates and returns a new Exporter which will collect metrics about zfs running on this machine
func New(poolNames []string) *Exporter {

	return &Exporter{
		poolNames: poolNames,
	}
}

// Describe sends all the descriptors of it's containted collectors to the provided channel
func (e *Exporter) Describe(c chan<- *prometheus.Desc) {

	z := NewZpool()
	z.Describe(c)
	d := NewDataset()
	d.Describe(c)
}

// Collect sends the collected metrics of all it's contained collectors on the provided channel
func (e *Exporter) Collect(c chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, name := range e.poolNames {

		// grab the pool
		pool, err := zfs.GetZpool(name)
		if err == nil {

			// grab the pools datasets
			ds, err := pool.Datasets()

			if err == nil {

				// metrics for the pool itself
				metrics := collectZpoolMetrics(pool, ds)
				for _, m := range metrics {
					c <- m
				}

				// metrics for each dataset in the pool
				for _, d := range ds {
					metrics = collectDatasetMetrics(d, pool)
					for _, m := range metrics {
						c <- m
					}
				}
			} else {
				log.Printf("Unable to collect datasets for zpool: %s %s\n", name, err)
			}
		} else {
			log.Printf("Unable to describe zpool: %s %s\n", name, err)
		}
	}
}
