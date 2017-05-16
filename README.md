# zfs_exporter

A Prometheus exporter for ZFS related high level metrics

## Why?

The regular [node_exporter](https://github.com/prometheus/node_exporter) from the [Prometheus project](https://prometheus.io) does include filesystem metrics as well as some ZFS related information.

But the ZFS information is very low level and not available on SmartOS or other Solaris-like systems.

The filesystem information especially does not reveal the important ZFS quota settings.

The original [zfs_exporter project](https://github.com/eliothedeman/zfs_exporter) was, however, not useable on some older ZFS tool stacks. This was mainly because the zpool tool used to have a bug: it claimed to support the "zpool get -p .." command but then complained about the "-p".

So we needed a "hack" to work around this bug until we could upgrade the affected operating systems.

## What?

The initial goals are:

- get the zfs_exporter working on older operating systems
- work around the "zpool get -p .." bug

## How?

The key steps so far have been

- fork all relevant Go projects
- include our hot fixes

In order to "activate" our fix you should run the tool via:

    ./zfs_exporter -feature.zpoolMetricsDisabled=true &
     
By default the exporter can be queried or tested with a request like this:

    curl http://localhost:9134/metrics
   
The TCP port 9134 is listed in the official [Prometheus port allocation](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) pages.
     
     
## Change Log

### v0.2.0

- Released on 2017.05.16
- Added the -feature.zpoolMetricsDisabled=true option.
- Added a filter to remove "snapshot" datasets from the metrics.
- Code now depends on a patched [go-zfs](https://github.com/tomi-engel/go-zfs) package 

Known Problems:

- Still does not work on 12.04 Ubuntu ZFS, because "/sbin/zfs zfs list -rHp ... => invalid option 'p'"


### v0.1.0

- Our initial fork from the original [zfs_exporter project](https://github.com/eliothedeman/zfs_exporter)


