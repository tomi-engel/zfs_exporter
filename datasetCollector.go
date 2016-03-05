package zfsexporter

import (
	"github.com/eliothedeman/go-zfs"
	"github.com/prometheus/client_golang/prometheus"
)

// A Dataset holds descriptions about a ZFS dataset for prometheus
type Dataset struct {
	Used          *prometheus.Desc
	Available     *prometheus.Desc
	Written       *prometheus.Desc
	VolumeSize    *prometheus.Desc
	UsedByDataset *prometheus.Desc
	LogicalUsed   *prometheus.Desc
	Quota         *prometheus.Desc
}

// Describe sends all of the descriptions of metrics collectd about datasets on the given channel
func (d *Dataset) Describe(c chan<- *prometheus.Desc) {
	m := []*prometheus.Desc{
		d.Used,
		d.Available,
		d.Written,
		d.VolumeSize,
		d.UsedByDataset,
		d.LogicalUsed,
		d.Quota,
	}

	for _, x := range m {
		c <- x
	}
}

// NewDataset fills in descriptions for a Dataset
func NewDataset() *Dataset {
	const (
		subsystem = "dataset"
	)

	labels := []string{
		"name",
		"pool_name",
		"type",
	}

	return &Dataset{
		Used: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"used"),
			"The amount of storage used by the underlying pool",
			labels,
			nil),
		Available: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"available"),
			"The amount of storage available to dataset",
			labels,
			nil),
		Written: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"written"),
			"The amount of storage written by this dataset",
			labels,
			nil),
		VolumeSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"volume_size"),
			"The size of the volume this datset resides on", // TODO figure out what this actually is
			labels,
			nil),
		UsedByDataset: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"used_by_datset"),
			"The amount of storage used by this dataset",
			labels,
			nil),
		LogicalUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"logical_used"),
			"The amount of storage used by this dataset (logically)",
			labels,
			nil),
		Quota: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,
				subsystem,
				"quota"),
			"The amount of storage the pool is allowed to use.", // figure out what this actually means
			labels,
			nil),
	}

}

func collectDatasetMetrics(ds *zfs.Dataset, pool *zfs.Zpool) []prometheus.Metric {
	desc := NewDataset()
	labels := []string{
		ds.Name,
		pool.Name,
		ds.Type,
	}

	return []prometheus.Metric{
		prometheus.MustNewConstMetric(
			desc.Used,
			prometheus.GaugeValue,
			float64(ds.Used),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Available,
			prometheus.GaugeValue,
			float64(ds.Avail),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Written,
			prometheus.CounterValue,
			float64(ds.Written),
			labels...),
		prometheus.MustNewConstMetric(
			desc.VolumeSize,
			prometheus.GaugeValue,
			float64(ds.Volsize),
			labels...),
		prometheus.MustNewConstMetric(
			desc.UsedByDataset,
			prometheus.GaugeValue,
			float64(ds.Usedbydataset),
			labels...),
		prometheus.MustNewConstMetric(
			desc.LogicalUsed,
			prometheus.GaugeValue,
			float64(ds.Logicalused),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Quota,
			prometheus.GaugeValue,
			float64(ds.Quota),
			labels...),
	}

}
