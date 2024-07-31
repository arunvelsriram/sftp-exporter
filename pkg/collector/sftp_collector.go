package collector

import (
	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
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

	objectCount = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "objects_available"),
		"Number of objects in the path",
		[]string{"path"},
		nil,
	)

	objectSize = prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, "", "objects_total_size_bytes"),
		"Total size of all the objects in the path",
		[]string{"path"},
		nil,
	)
)

type SFTPCollector struct {
	sftpClient client.SFTPClient
}

func (s SFTPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	useStatVfs := viper.GetBool(viperkeys.SFTPStatVfs)
	if useStatVfs {
		ch <- fsTotalSpace
		ch <- fsFreeSpace
	}
	ch <- objectCount
	ch <- objectSize
}

func (s SFTPCollector) Collect(ch chan<- prometheus.Metric) {
	paths := viper.GetStringSlice(viperkeys.SFTPPaths)

	if err := s.sftpClient.Connect(); err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		log.WithField("when", "collecting up metric").Error(err)
		return
	}
	defer s.sftpClient.Close()
	log.Debug("connected to SFTP")
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	useStatVfs := viper.GetBool(viperkeys.SFTPStatVfs)
	if useStatVfs {
		log.Debug("collecting filesystem metrics")
		for _, path := range paths {
			log.Debugf("collecting filesystem metrics for path: %s", path)
			statVFS, err := s.sftpClient.StatVFS(path)
			if err != nil {
				log.WithFields(log.Fields{"when": "collecting filesystem metrics", "path": path}).Error(err)
			} else {
				totalSpace := float64(statVFS.TotalSpace())
				freeSpace := float64(statVFS.FreeSpace())
				log.Debugf("writing filesystem metrics for path: %s", path)
				ch <- prometheus.MustNewConstMetric(fsTotalSpace, prometheus.GaugeValue, totalSpace, path)
				ch <- prometheus.MustNewConstMetric(fsFreeSpace, prometheus.GaugeValue, freeSpace, path)
			}
		}
	}

	log.Debug("collecting object metrics")
	for _, path := range paths {
		log.Debugf("collecting object metrics for path: %s", path)
		var size int64
		count := 0
		var walkErr error
		walker := s.sftpClient.Walk(path)
		for walker.Step() {
			if walkErr = walker.Err(); walkErr != nil {
				log.WithFields(log.Fields{"when": "collecting object metrics", "path": path}).Error(walkErr)
				break
			}

			if walker.Stat().IsDir() {
				continue
			}
			size += walker.Stat().Size()
			count++
		}
		if walkErr == nil {
			ch <- prometheus.MustNewConstMetric(objectCount, prometheus.GaugeValue, float64(count), path)
			ch <- prometheus.MustNewConstMetric(objectSize, prometheus.GaugeValue, float64(size), path)
		}
	}
}

func NewSFTPCollector(c client.SFTPClient) prometheus.Collector {
	return SFTPCollector{sftpClient: c}
}
