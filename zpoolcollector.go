package zfsexporter

import (
	"github.com/eliothedeman/go-zfs"
	"github.com/prometheus/client_golang/prometheus"
)

func describePool(pool *zfs.Zpool) map[string]*prometheus.Desc {
	const (
		subsystem = "zpool"
	)

	labels := []string{
		"pool_name",
		"hostname",
	}

	return map[string]*prometheus.Desc{
		"allocated": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "allocated"),
			"Bytes of storage physically allocated",
			labels,
			nil),
		"size": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "size"),
			"Total Size of the storage pool",
			labels,
			nil),
		"free": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "free"),
			"Bytes of storage free in this pool",
			labels,
			nil),
		"fragmentation": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "fragmentation"),
			"Amount of fragmentation in a pool", // TODO figure out what the metric actually means for this
			labels,
			nil),
		"readonly": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "readonly"),
			"True if pool is in readonly mode",
			labels,
			nil),
		"freeing": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "freeing"),
			"Amount of storage currently being freed", // TODO figure out what the metric actually means for this
			labels,
			nil),
		"leaked": prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "leaked"),
			"", // TODO figure out what the metric actually means for this
			labels,
			nil),
		"dedup_ratio": prometheus.NewDesc(
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
			desc["allocated"],
			prometheus.GaugeValue,
			float64(pool.Allocated),
			labels...),
		prometheus.MustNewConstMetric(
			desc["size"],
			prometheus.GaugeValue,
			float64(pool.Size),
			labels...),
		prometheus.MustNewConstMetric(
			desc["free"],
			prometheus.GaugeValue,
			float64(pool.Free),
			labels...),
		prometheus.MustNewConstMetric(
			desc["fragmentation"],
			prometheus.GaugeValue,
			float64(pool.Fragmentation),
			labels...),
		prometheus.MustNewConstMetric(
			desc["freeing"],
			prometheus.GaugeValue,
			float64(pool.Freeing),
			labels...),
		prometheus.MustNewConstMetric(
			desc["leaked"],
			prometheus.GaugeValue,
			float64(pool.Leaked),
			labels...),
		prometheus.MustNewConstMetric(
			desc["dedup_ratio"],
			prometheus.GaugeValue,
			float64(pool.DedupRatio),
			labels...),
	}
}
