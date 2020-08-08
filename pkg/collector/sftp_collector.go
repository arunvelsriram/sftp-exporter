package collector

import (
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	buildInfo = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "build_info"),
		"",
		[]string{"version"},
		nil,
	)
)

type SFTPCollector struct {
}

func (c SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- buildInfo
}

func (c SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(buildInfo, prometheus.GaugeValue, 1, "dummy")
}

func NewSFTPCollector() prometheus.Collector {
	return SFTPCollector{}
}
