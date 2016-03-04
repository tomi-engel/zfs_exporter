package zfsexporter

import (
	"github.com/eliothedeman/go-zfs"
	"github.com/prometheus/client_golang/prometheus"
)

// A Zpool holds descriptions about metrics about a zpool
type Zpool struct {
	Allocated     *prometheus.Desc
	Size          *prometheus.Desc
	Free          *prometheus.Desc
	Fragmentation *prometheus.Desc
	ReadOnly      *prometheus.Desc
	Freeing       *prometheus.Desc
	Leaked        *prometheus.Desc
	DedupRatio    *prometheus.Desc
}

func describePool(pool *zfs.Zpool) Zpool {
	const (
		subsystem = "zpool"
	)

	labels := []string{
		"pool_name",
		"hostname",
	}

	return Zpool{
		Allocated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "allocated"),
			"Bytes of storage physically allocated",
			labels,
			nil),
		Size: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "size"),
			"Total Size of the storage pool",
			labels,
			nil),
		Free: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "free"),
			"Bytes of storage free in this pool",
			labels,
			nil),
		Fragmentation: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "fragmentation"),
			"Amount of fragmentation in a pool", // TODO figure out what the metric actually means for this
			labels,
			nil),
		ReadOnly: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "readonly"),
			"True if pool is in readonly mode",
			labels,
			nil),
		Freeing: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "freeing"),
			"Amount of storage currently being freed", // TODO figure out what the metric actually means for this
			labels,
			nil),
		Leaked: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "leaked"),
			"", // TODO figure out what the metric actually means for this
			labels,
			nil),
		DedupRatio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "dedup_ratio"),
			"", // TODO figure out what the metric actually means for this
			labels,
			nil),
	}
}

func createMetrics(pool *zfs.Zpool, hostname string) []prometheus.Metric {

	desc := describePool(pool)
	labels := []string{
		pool.Name,
		hostname,
	}
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(
			desc.Allocated,
			prometheus.GaugeValue,
			float64(pool.Allocated),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Size,
			prometheus.GaugeValue,
			float64(pool.Size),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Size,
			prometheus.GaugeValue,
			float64(pool.Free),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Fragmentation,
			prometheus.GaugeValue,
			float64(pool.Fragmentation),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Freeing,
			prometheus.GaugeValue,
			float64(pool.Freeing),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Leaked,
			prometheus.GaugeValue,
			float64(pool.Leaked),
			labels...),
		prometheus.MustNewConstMetric(
			desc.DedupRatio,
			prometheus.GaugeValue,
			float64(pool.DedupRatio),
			labels...),
	}
}
