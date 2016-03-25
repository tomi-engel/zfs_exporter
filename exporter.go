package zfsexporter

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "zfs"
)

// An Exporter is a Prometheus exporter for the ZFS filesystem.
type Exporter struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

// Make sure the exporter satisfies the prometheus collector interface
var _ prometheus.Collector = &Exporter{}

// New creates and returns a new Exporter which will collect metrics
// about ZFS zpools and datasets running on this machine.
func New(pools []string) *Exporter {
	return &Exporter{
		collectors: []prometheus.Collector{
			NewZpoolCollector(pools),
			NewDatasetCollector(pools),
		},
	}
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, cc := range e.collectors {
		cc.Describe(ch)
	}
}

// Collect sends the collected metrics from each of the collectors to
// prometheus. Collect could be called several times concurrently
// and thus its run is protected by a single mutex.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, cc := range e.collectors {
		cc.Collect(ch)
	}
}
