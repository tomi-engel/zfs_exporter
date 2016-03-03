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

// make sure the exporter satisfies the prometheus collector interface
var _ prometheus.Collector = &Exporter{}

// New creates and returns a new Exporter which will collect metrics about zfs running on this machine
func New() *Exporter {
	return &Exporter{
		collectors: []prometheus.Collector{},
	}
}

// Describe sends all the descriptors of it's containted collectors to the provided channel
func (e *Exporter) Describe(c chan<- *prometheus.Desc) {

	// lock because we can't trust that collectors won't be added while we describe
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, col := range e.collectors {
		col.Describe(c)
	}
}

// Collect sends the collected metrics of all it's contained collectors on the provided channel
func (e *Exporter) Collect(c chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, col := range e.collectors {
		col.Collect(c)
	}
}
