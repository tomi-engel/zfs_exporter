package zfsexporter

import (
	"github.com/eliothedeman/go-zfs"
	"github.com/prometheus/client_golang/prometheus"
)

// A Filesystem holds descriptions about a ZFS dataset for prometheus
type Filesystem struct {
	Used          *prometheus.Desc
	Available     *prometheus.Desc
	Written       *prometheus.Desc
	VolumeSize    *prometheus.Desc
	UsedByDataset *prometheus.Desc
	LogicalUsed   *prometheus.Desc
	Quota         *prometheus.Desc
}

func describeFileSystem(dataset *zfs.Dataset) Filesystem {
	const (
		subsystem = "filesystem"
	)

	labels := []string{
		"name",
		"origin",
		"pool_name",
		"hostname",
		"mount_point",
	}

	return Filesystem{
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
			"", // TODO figure out what this actually is
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
			"", // figure out what this actually means
			labels,
			nil),
	}

}

func createFilesystemMetrics(fs *zfs.Dataset, pool *zfs.Zpool, hostname string) []prometheus.Metric {
	desc := describeFileSystem(fs)
	labels := []string{
		fs.Name,
		fs.Origin,
		pool.Name,
		hostname,
		fs.Mountpoint,
	}

	return []prometheus.Metric{
		prometheus.MustNewConstMetric(
			desc.Used,
			prometheus.GaugeValue,
			float64(fs.Used),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Available,
			prometheus.GaugeValue,
			float64(fs.Avail),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Written,
			prometheus.GaugeValue,
			float64(fs.Written),
			labels...),
		prometheus.MustNewConstMetric(
			desc.VolumeSize,
			prometheus.GaugeValue,
			float64(fs.Volsize),
			labels...),
		prometheus.MustNewConstMetric(
			desc.UsedByDataset,
			prometheus.GaugeValue,
			float64(fs.Usedbydataset),
			labels...),
		prometheus.MustNewConstMetric(
			desc.LogicalUsed,
			prometheus.GaugeValue,
			float64(fs.Logicalused),
			labels...),
		prometheus.MustNewConstMetric(
			desc.Quota,
			prometheus.GaugeValue,
			float64(fs.Quota),
			labels...),
	}

}
