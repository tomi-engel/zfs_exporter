package zfsexporter

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tomi-engel/go-zfs"
)

// A ZpoolCollector is a Prometheus collector for ZFS zpool metrics.
type ZpoolCollector struct {
	AllocatedBytes       *prometheus.Desc
	SizeBytes            *prometheus.Desc
	FreeBytes            *prometheus.Desc
	FragmentationPercent *prometheus.Desc
	ReadOnly             *prometheus.Desc
	FreeingBytes         *prometheus.Desc
	LeakedBytes          *prometheus.Desc
	DeduplicationRatio   *prometheus.Desc
	Snapshots            *prometheus.Desc
	Filesystems          *prometheus.Desc
	Volumes              *prometheus.Desc

	pools []string
}

// NewZpoolCollector creates a new ZpoolCollector.
func NewZpoolCollector(pools []string) *ZpoolCollector {
	const (
		subsystem = "zpool"
	)

	labels := []string{
		"pool",
	}

	return &ZpoolCollector{
		AllocatedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "allocated_bytes"),
			"Number of allocated bytes in the zpool",
			labels,
			nil,
		),

		SizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "size_bytes"),
			"Number of total bytes in the zpool",
			labels,
			nil,
		),

		FreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "free_bytes"),
			"Number of free bytes in the zpool",
			labels,
			nil,
		),

		FragmentationPercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "fragmentation_percent"),
			"Fragmentation percentage for the zpool",
			labels,
			nil,
		),

		ReadOnly: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "readonly"),
			"Whether or not the zpool is read-only; 1 if it is read-only, 0 otherwise",
			labels,
			nil,
		),

		FreeingBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "freeing_bytes"),
			"Number of bytes currently being freed in the zpool",
			labels,
			nil,
		),

		LeakedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "leaked_bytes"),
			"Number of bytes of leaked storage in the zpool",
			labels,
			nil,
		),

		DeduplicationRatio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "deduplication_ratio"),
			"Ratio of deduplicated content in the zpool",
			labels,
			nil,
		),

		Snapshots: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "snapshots"),
			"Total number of snapshots in the zpool",
			labels,
			nil,
		),

		Filesystems: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "filesystems"),
			"Total number of filesystems in the zpool",
			labels,
			nil,
		),

		Volumes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "volumes"),
			"Total number of volumes in the zpool",
			labels,
			nil,
		),

		pools: pools,
	}
}

// collect begins a metrics collection task for all metrics related to UniFi
// stations.
func (c *ZpoolCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	for _, p := range c.pools {
		zpool, err := zfs.GetZpool(p)
		if err != nil {
			return c.SizeBytes, err
		}

		ds, err := zpool.Datasets()
		if err != nil {
			return c.Snapshots, err
		}

		c.collectZpoolMetrics(ch, zpool, ds)
	}

	return nil, nil
}

// collectZpoolMetrics collects metrics for an individual zpool.
func (c *ZpoolCollector) collectZpoolMetrics(ch chan<- prometheus.Metric, zpool *zfs.Zpool, ds []*zfs.Dataset) {
	labels := []string{
		zpool.Name,
	}

	ch <- prometheus.MustNewConstMetric(
		c.AllocatedBytes,
		prometheus.GaugeValue,
		float64(zpool.Allocated),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.SizeBytes,
		prometheus.GaugeValue,
		float64(zpool.Size),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FreeBytes,
		prometheus.GaugeValue,
		float64(zpool.Free),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FragmentationPercent,
		prometheus.GaugeValue,
		float64(zpool.Fragmentation),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FreeingBytes,
		prometheus.GaugeValue,
		float64(zpool.Freeing),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.LeakedBytes,
		prometheus.GaugeValue,
		float64(zpool.Leaked),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.DeduplicationRatio,
		prometheus.GaugeValue,
		float64(zpool.DedupRatio),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Snapshots,
		prometheus.GaugeValue,
		float64(countDatasetsByType(ds, zfs.DatasetSnapshot)),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Filesystems,
		prometheus.GaugeValue,
		float64(countDatasetsByType(ds, zfs.DatasetFilesystem)),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.Volumes,
		prometheus.GaugeValue,
		float64(countDatasetsByType(ds, zfs.DatasetVolume)),
		labels...,
	)
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *ZpoolCollector) Describe(ch chan<- *prometheus.Desc) {
	m := []*prometheus.Desc{
		c.AllocatedBytes,
		c.SizeBytes,
		c.FreeBytes,
		c.FragmentationPercent,
		c.ReadOnly,
		c.FreeingBytes,
		c.LeakedBytes,
		c.DeduplicationRatio,
		c.Snapshots,
		c.Filesystems,
		c.Volumes,
	}

	for _, d := range m {
		ch <- d
	}
}

// Collect sends the metric values for each metric pertaining to ZFS zpools
// over to the provided prometheus Metric channel.
func (c *ZpoolCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting zpool metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}

// countDatasetsByType retrieves a count of datasets which match the
// specified type.
func countDatasetsByType(ds []*zfs.Dataset, dsType string) int {
	var count int
	for _, d := range ds {
		if d.Type == dsType {
			count++
		}
	}

	return count
}
