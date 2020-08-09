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
		[]string{},
		nil,
	)

	fsFreeSpace = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "filesystem_free_space_bytes"),
		"Free space in the filesystem containing the path",
		[]string{},
		nil,
	)
)

type SFTPCollector struct {
	config        config.Config
	clientFactory client.Factory
}

func (s SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (s SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	sftpClient, err := s.clientFactory.SFTPClient()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		log.WithFields(log.Fields{"event": "creating SFTP sftpClient"}).Error(err)
		return
	}
	defer sftpClient.Close()
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	fsStat, err := sftpClient.FSStat()
	if err != nil {
		log.WithFields(log.Fields{"event": "getting FS stat"}).Error(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(fsTotalSpace, prometheus.GaugeValue, fsStat.TotalSpace)
	ch <- prometheus.MustNewConstMetric(fsFreeSpace, prometheus.GaugeValue, fsStat.FreeSpace)
}

func NewSFTPCollector(cfg config.Config, f client.Factory) prometheus.Collector {
	return SFTPCollector{
		config:        cfg,
		clientFactory: f,
	}
}
