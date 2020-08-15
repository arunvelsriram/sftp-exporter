package collector

import (
	"github.com/arunvelsriram/sftp-exporter/pkg/service"
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
		prometheus.BuildFQName(c.Namespace, "", "objects_size_total_bytes"),
		"Total size of all objects in the path",
		[]string{"path"},
		nil,
	)
)

type CreateClientFn func(config.Config) (client.SFTPClient, error)

type SFTPCollector struct {
	config      config.Config
	sftpService service.SFTPService
}

func (s SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- fsTotalSpace
	ch <- fsFreeSpace
	ch <- objectCount
	ch <- objectSize
}

func (s SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	if err := s.sftpService.Connect(); err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		log.WithFields(log.Fields{"event": "creating SFTP connection"}).Error(err)
		return
	}
	defer s.sftpService.Close()
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	fsStats := s.sftpService.FSStats()
	for _, stat := range fsStats {
		ch <- prometheus.MustNewConstMetric(fsTotalSpace, prometheus.GaugeValue, stat.TotalSpace, stat.Path)
		ch <- prometheus.MustNewConstMetric(fsFreeSpace, prometheus.GaugeValue, stat.FreeSpace, stat.Path)
	}

	objectStats := s.sftpService.ObjectStats()
	for _, stat := range objectStats {
		ch <- prometheus.MustNewConstMetric(objectCount, prometheus.GaugeValue, stat.ObjectCount, stat.Path)
		ch <- prometheus.MustNewConstMetric(objectSize, prometheus.GaugeValue, stat.ObjectSize, stat.Path)
	}
}

func NewSFTPCollector(cfg config.Config, s service.SFTPService) prometheus.Collector {
	return SFTPCollector{
		config:      cfg,
		sftpService: s,
	}
}
