package collector

import (
	"fmt"
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
	"github.com/arunvelsriram/sftp-exporter/pkg/model"
	"github.com/pkg/sftp"
	"github.com/spf13/viper"

	"github.com/arunvelsriram/sftp-exporter/pkg/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/suite"

	log "github.com/sirupsen/logrus"
)

type SFTPCollectorSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	sftpService *mocks.MockSFTPService
	sftpClient  *mocks.MockSFTPClient
	collector   prometheus.Collector
}

func TestSFTPCollectorSuite(t *testing.T) {
	suite.Run(t, new(SFTPCollectorSuite))
}

func (s *SFTPCollectorSuite) SetupTest() {
	log.SetLevel(log.DebugLevel)
	s.ctrl = gomock.NewController(s.T())
	s.sftpService = mocks.NewMockSFTPService(s.ctrl)
	s.sftpClient = mocks.NewMockSFTPClient(s.ctrl)
	s.collector = NewSFTPCollector(s.sftpService, s.sftpClient)
}

func (s *SFTPCollectorSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SFTPCollectorSuite) TestSFTPCollectorDescribe() {
	ch := make(chan *prometheus.Desc)
	go s.collector.Describe(ch)

	up := <-ch
	s.Equal(`Desc{fqName: "sftp_up", help: "Tells if exporter is able to connect to SFTP", `+
		`constLabels: {}, variableLabels: []}`,
		up.String(),
	)

	fsTotalSpace := <-ch
	s.Equal(`Desc{fqName: "sftp_filesystem_total_space_bytes", `+
		`help: "Total space in the filesystem containing the path", constLabels: {}, variableLabels: [path]}`,
		fsTotalSpace.String(),
	)

	fsFreeSpace := <-ch
	s.Equal(`Desc{fqName: "sftp_filesystem_free_space_bytes", `+
		`help: "Free space in the filesystem containing the path", constLabels: {}, variableLabels: [path]}`,
		fsFreeSpace.String(),
	)

	objectCount := <-ch
	s.Equal(
		`Desc{fqName: "sftp_objects_count_total", `+
			`help: "Total number of objects in the path", constLabels: {}, variableLabels: [path]}`,
		objectCount.String(),
	)

	objectSize := <-ch
	s.Equal(
		`Desc{fqName: "sftp_objects_size_total_bytes", `+
			`help: "Total size of all objects in the path", constLabels: {}, variableLabels: [path]}`,
		objectSize.String(),
	)
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldWriteUpMetric() {
	viper.Set(viperkeys.SFTPPaths, []string{})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpService.EXPECT().ObjectStats()
	s.sftpClient.EXPECT().Close().Return(nil)
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	up := <-ch
	metric := dto.Metric{}
	desc := up.Desc()
	_ = up.Write(&metric)
	s.Equal(`Desc{fqName: "sftp_up", help: "Tells if exporter is able to connect to SFTP", `+
		`constLabels: {}, variableLabels: []}`, desc.String())
	s.Equal(1.0, metric.GetGauge().GetValue())

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldWriteUpMetricAndReturnIfClientCreationFails() {
	viper.Set(viperkeys.SFTPPaths, []string{})
	s.sftpClient.EXPECT().Connect().Return(fmt.Errorf("failed to connect to SFTP"))
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	up := <-ch
	metric := dto.Metric{}
	desc := up.Desc()
	_ = up.Write(&metric)
	s.Equal(`Desc{fqName: "sftp_up", help: "Tells if exporter is able to connect to SFTP", `+
		`constLabels: {}, variableLabels: []}`, desc.String())
	s.Equal(0.0, metric.GetGauge().GetValue())

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldWriteFSMetrics() {
	viper.Set(viperkeys.SFTPPaths, []string{"/path0", "/path1"})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpClient.EXPECT().StatVFS("/path0").Return(&sftp.StatVFS{Frsize: 10, Blocks: 1000, Bfree: 100}, nil)
	s.sftpClient.EXPECT().StatVFS("/path1").Return(&sftp.StatVFS{Frsize: 5, Blocks: 1000, Bfree: 500}, nil)
	s.sftpService.EXPECT().ObjectStats()
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	<-ch
	metric := &dto.Metric{}
	var desc *prometheus.Desc

	totalSpace1 := <-ch
	desc = totalSpace1.Desc()
	_ = totalSpace1.Write(metric)
	s.Equal(`Desc{fqName: "sftp_filesystem_total_space_bytes", help: "Total space in the filesystem containing the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(10000.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path0", metric.GetLabel()[0].GetValue())

	freeSpace1 := <-ch
	desc = freeSpace1.Desc()
	_ = freeSpace1.Write(metric)
	s.Equal(`Desc{fqName: "sftp_filesystem_free_space_bytes", help: "Free space in the filesystem containing the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(1000.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path0", metric.GetLabel()[0].GetValue())

	totalSpace2 := <-ch
	desc = totalSpace2.Desc()
	_ = totalSpace2.Write(metric)
	s.Equal(`Desc{fqName: "sftp_filesystem_total_space_bytes", help: "Total space in the filesystem containing the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(5000.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path1", metric.GetLabel()[0].GetValue())

	freeSpace2 := <-ch
	desc = freeSpace2.Desc()
	_ = freeSpace2.Write(metric)
	s.Equal(`Desc{fqName: "sftp_filesystem_free_space_bytes", help: "Free space in the filesystem containing the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(2500.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path1", metric.GetLabel()[0].GetValue())

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldNotWriteFSMetricsForError() {
	viper.Set(viperkeys.SFTPPaths, []string{"/path0"})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpClient.EXPECT().StatVFS("/path0").Return(nil, fmt.Errorf("failed to get VFS stats"))
	s.sftpService.EXPECT().ObjectStats()
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	<-ch
	close(ch)

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldWriteObjectStats() {
	viper.Set(viperkeys.SFTPPaths, []string{})
	objectStats := model.ObjectStats([]model.ObjectStat{
		{
			Path:        "/path1",
			ObjectCount: 100,
			ObjectSize:  1111.11,
		},
		{
			Path:        "/path2",
			ObjectCount: 200,
			ObjectSize:  2222.22,
		},
	})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpService.EXPECT().ObjectStats().Return(objectStats)
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	<-ch
	metric := &dto.Metric{}
	var desc *prometheus.Desc

	objectCount1 := <-ch
	desc = objectCount1.Desc()
	_ = objectCount1.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_count_total", help: "Total number of objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(100.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path1", metric.GetLabel()[0].GetValue())

	objectSize1 := <-ch
	desc = objectSize1.Desc()
	_ = objectSize1.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_size_total_bytes", help: "Total size of all objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(1111.11, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path1", metric.GetLabel()[0].GetValue())

	objectCount2 := <-ch
	desc = objectCount2.Desc()
	_ = objectCount2.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_count_total", help: "Total number of objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(200.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path2", metric.GetLabel()[0].GetValue())

	objectSize2 := <-ch
	desc = objectSize2.Desc()
	_ = objectSize2.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_size_total_bytes", help: "Total size of all objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(2222.22, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path2", metric.GetLabel()[0].GetValue())

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldNotWriteObjectStatsForEmpty() {
	viper.Set(viperkeys.SFTPPaths, []string{})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpService.EXPECT().ObjectStats().Return(model.ObjectStats(model.ObjectStats{}))
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	<-ch
	close(ch)

	<-done
}
