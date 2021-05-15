package collector

import (
	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
	"github.com/arunvelsriram/sftp-exporter/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/arunvelsriram/sftp-exporter/pkg/client"
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

	// TODO: Rename this metric and re-evaluate implementation
	objectCount = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "objects_count_total"),
		"Total number of objects in the path",
		[]string{"path"},
		nil,
	)

	// TODO: Rename this metric and re-evaluate implementation
	objectSize = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "objects_size_total_bytes"),
		"Total size of all objects in the path",
		[]string{"path"},
		nil,
	)
)

type SFTPCollector struct {
	sftpService service.SFTPService
	sftpClient  client.SFTPClient
}

func (s SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- fsTotalSpace
	ch <- fsFreeSpace
	ch <- objectCount
	ch <- objectSize
}

func (s SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	if err := s.sftpClient.Connect(); err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		log.WithField("when", "collecting up metric").Error(err)
		return
	}
	defer s.sftpClient.Close()
	log.Debug("connected to SFTP")
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	log.Debug("collecting filesystem metrics")
	paths := viper.GetStringSlice(viperkeys.SFTPPaths)
	for _, path := range paths {
		log.Debugf("collecting filesystem metrics for path: %s", path)
		statVFS, err := s.sftpClient.StatVFS(path)
		if err != nil {
			log.WithField("when", "collecting filesystem metrics").Error(err)
		} else {
			totalSpace := float64(statVFS.TotalSpace())
			freeSpace := float64(statVFS.FreeSpace())
			log.Debugf("writing filesystem metrics for path: %s", path)
			ch <- prometheus.MustNewConstMetric(fsTotalSpace, prometheus.GaugeValue, totalSpace, path)
			ch <- prometheus.MustNewConstMetric(fsFreeSpace, prometheus.GaugeValue, freeSpace, path)
		}
	}

	log.Debug("collecting object metrics")
	objectStats := s.sftpService.ObjectStats()
	for _, stat := range objectStats {
		ch <- prometheus.MustNewConstMetric(objectCount, prometheus.GaugeValue, stat.ObjectCount, stat.Path)
		ch <- prometheus.MustNewConstMetric(objectSize, prometheus.GaugeValue, stat.ObjectSize, stat.Path)
	}
}

func NewSFTPCollector(s service.SFTPService, c client.SFTPClient) prometheus.Collector {
	return SFTPCollector{sftpService: s, sftpClient: c}
}
