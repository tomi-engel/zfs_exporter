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
	Snapshots     *prometheus.Desc
	Filesystems   *prometheus.Desc
	Volumes       *prometheus.Desc
}

// Describe sends the descriptions of the zpool on the given channel
func (z *Zpool) Describe(c chan<- *prometheus.Desc) {
	m := []*prometheus.Desc{
		z.Allocated,
		z.Size,
		z.Free,
		z.Fragmentation,
		z.ReadOnly,
		z.Freeing,
		z.Leaked,
		z.DedupRatio,
		z.Snapshots,
		z.Filesystems,
		z.Volumes,
	}
	for _, d := range m {
		c <- d
	}
}

// NewZpool fills a zpool with it's descriptions
func NewZpool() *Zpool {
	const (
		subsystem = "zpool"
	)

	labels := []string{
		"pool_name",
	}

	return &Zpool{
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
			"Amount of leaked storage", // TODO figure out what the metric actually means for this
			labels,
			nil),
		DedupRatio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "dedup_ratio"),
			"Ratio of storage that is used for duplication", // TODO figure out what the metric actually means for this
			labels,
			nil),
		Snapshots: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "snapshots"),
			"The number of snapshots in this pool",
			labels,
			nil),
		Filesystems: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "filesystems"),
			"The number of filesystems in this pool",
			labels,
			nil),
		Volumes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "volumes"),
			"The number of volumes in this pool",
			labels,
			nil),
	}
}

func collectZpoolMetrics(pool *zfs.Zpool, ds []*zfs.Dataset) []prometheus.Metric {

	desc := NewZpool()
	labels := []string{
		pool.Name,
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
		prometheus.MustNewConstMetric(
			desc.Snapshots,
			prometheus.GaugeValue,
			float64(countDatasetsByType(ds, zfs.DatasetSnapshot)),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Filesystems,
			prometheus.GaugeValue,
			float64(countDatasetsByType(ds, zfs.DatasetFilesystem)),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Volumes,
			prometheus.GaugeValue,
			float64(countDatasetsByType(ds, zfs.DatasetVolume)),
			labels...),
	}
}

func countDatasetsByType(ds []*zfs.Dataset, dsType string) int {
	count := 0
	for _, d := range ds {
		if d.Type == dsType {
			count++
		}
	}
	return count
}
