package collector

import (
	"fmt"

	"github.com/arunvelsriram/sftp-exporter/pkg/client"
	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "up"),
		"Tells if exporter is able to connect to SFTP",
		[]string{},
		nil,
	)

	fsTotalSpace = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "filesystem_total_space_bytes"),
		"Total space in the filesystem containing user directory",
		[]string{},
		nil,
	)

	fsFreeSpace = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "filesystem_free_space_bytes"),
		"Free space in the filesystem containing user directory",
		[]string{},
		nil,
	)
)

type SFTPCollector struct {
	config config.Config
}

func (s SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (s SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	client, err := client.NewSFTPClient(s.config)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		fmt.Println("failed to get create sftp connection")
		fmt.Println(err)
		client.Close()
		return
	}
	defer client.Close()
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	fsStat, err := client.FSStat()
	if err != nil {
		fmt.Println("failed to get FS stat")
		fmt.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(fsTotalSpace, prometheus.GaugeValue, fsStat.TotalSpace)
	ch <- prometheus.MustNewConstMetric(fsFreeSpace, prometheus.GaugeValue, fsStat.FreeSpace)
}

func NewSFTPCollector(cfg config.Config) prometheus.Collector {
	return SFTPCollector{cfg}
}
