package collector

import (
	log "github.com/sirupsen/logrus"

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
		"Total space in the filesystem containing the path",
		[]string{"path"},
		nil,
	)

	fsFreeSpace = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "filesystem_free_space_bytes"),
		"Free space in the filesystem containing the path",
		[]string{"path"},
		nil,
	)

	objectCount = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "objects_count_total"),
		"Total number of objects in the path",
		[]string{"path"},
		nil,
	)

	objectSize = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "objects_size_total"),
		"Total size of all objects in the path",
		[]string{"path"},
		nil,
	)
)

type SFTPCollector struct {
	config config.Config
}

func (s SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- fsTotalSpace
	ch <- fsFreeSpace
	ch <- objectCount
	ch <- objectSize
}

func (s SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	sftpClient, err := client.NewSFTPClient(s.config)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		log.WithFields(log.Fields{"event": "creating SFTP sftpClient"}).Error(err)
		return
	}
	defer sftpClient.Close()
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	stats, err := sftpClient.FSStat()
	if err != nil {
		log.WithFields(log.Fields{"event": "getting FS stats"}).Error(err)
		return
	}
	for _, stat := range stats {
		ch <- prometheus.MustNewConstMetric(fsTotalSpace, prometheus.GaugeValue, stat.TotalSpace, stat.Path)
		ch <- prometheus.MustNewConstMetric(fsFreeSpace, prometheus.GaugeValue, stat.FreeSpace, stat.Path)
	}
}

func NewSFTPCollector(cfg config.Config) prometheus.Collector {
	return SFTPCollector{config: cfg}
}
