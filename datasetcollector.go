package zfsexporter

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tomi-engel/go-zfs"
)

// A DatasetCollector is a Prometheus collector for ZFS dataset metrics.
type DatasetCollector struct {
	UsedBytes        *prometheus.Desc
	AvailableBytes   *prometheus.Desc
	WrittenBytes     *prometheus.Desc
	LogicalUsedBytes *prometheus.Desc
	QuotaBytes       *prometheus.Desc

	pools []string
}

// NewDatasetCollector creates a new DatasetCollector.
func NewDatasetCollector(pools []string) *DatasetCollector {
	const (
		subsystem = "dataset"
	)

	labels := []string{
		"name",
		"pool",
		"type",
	}

	return &DatasetCollector{
		UsedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "used_bytes"),
			"Number of used bytes in the dataset",
			labels,
			nil,
		),

		AvailableBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "available_bytes"),
			"Number of bytes available to the dataset",
			labels,
			nil,
		),

		WrittenBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "written_bytes"),
			"Number of bytes written to the dataset",
			labels,
			nil,
		),

		LogicalUsedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "logical_used_bytes"),
			"Number of logically used bytes in the dataset",
			labels,
			nil,
		),

		QuotaBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "quota_bytes"),
			"Number of quota bytes available to the dataset",
			labels,
			nil,
		),

		pools: pools,
	}
}

// collect begins a metrics collection task for all metrics related to ZFS
// datasets.
func (c *DatasetCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	for _, p := range c.pools {

		var zpool *zfs.Zpool
		var err error
		zpool, err = nil, nil

		if *featureZpoolMetricsDisabled == true {
			zpool, err = zfs.GetZpoolWithoutUsingGetPFeature(p)
		} else {
			zpool, err = zfs.GetZpool(p)
		}
		if err != nil {
			return c.UsedBytes, err
		}

		ds, err := zpool.Datasets()
		if err != nil {
			return c.UsedBytes, err
		}

		for _, d := range ds {

			// We have way too many snapshots which are changing way too quickly
			// this is "bad" for the Prometheus data storage subsystem.
			// Besides that .. snapshots are not changing, so tracing them has little value.
			// So lets exclude the snapshots..

			if d.Type != "snapshot" {
				c.collectDatasetMetrics(ch, zpool, d)
			}
		}
	}

	return nil, nil
}

// collectDatasetMetrics collects metrics for an individual dataset.
func (c *DatasetCollector) collectDatasetMetrics(ch chan<- prometheus.Metric, zpool *zfs.Zpool, ds *zfs.Dataset) {
	labels := []string{
		ds.Name,
		zpool.Name,
		ds.Type,
	}

	ch <- prometheus.MustNewConstMetric(
		c.UsedBytes,
		prometheus.GaugeValue,
		float64(ds.Used),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.AvailableBytes,
		prometheus.GaugeValue,
		float64(ds.Avail),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.WrittenBytes,
		prometheus.CounterValue,
		float64(ds.Written),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.LogicalUsedBytes,
		prometheus.GaugeValue,
		float64(ds.Logicalused),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.QuotaBytes,
		prometheus.GaugeValue,
		float64(ds.Quota),
		labels...,
	)
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *DatasetCollector) Describe(ch chan<- *prometheus.Desc) {
	m := []*prometheus.Desc{
		c.UsedBytes,
		c.AvailableBytes,
		c.WrittenBytes,
		c.LogicalUsedBytes,
		c.QuotaBytes,
	}

	for _, x := range m {
		ch <- x
	}
}

// Collect sends the metric values for each metric pertaining to ZFS datasets
// over to the provided prometheus Metric channel.
func (c *DatasetCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting dataset metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}
