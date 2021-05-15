package collector

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
	"github.com/arunvelsriram/sftp-exporter/pkg/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/kr/fs"
	"github.com/pkg/sftp"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	log "github.com/sirupsen/logrus"
)

// Refer: https://github.com/kr/fs/blob/main/filesystem.go
type memKrFs struct {
	memFs afero.Fs
}

func (m memKrFs) ReadDir(dirname string) ([]os.FileInfo, error) {
	if strings.EqualFold("/errorpath", dirname) {
		return nil, fmt.Errorf("error reading directory")
	}
	return afero.ReadDir(m.memFs, dirname)
}

func (m memKrFs) Lstat(name string) (os.FileInfo, error) {
	return m.memFs.Stat(name)
}

func (m memKrFs) Join(elem ...string) string {
	return path.Join(elem...)
}

type SFTPCollectorSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	sftpClient *mocks.MockSFTPClient
	collector  prometheus.Collector
}

func TestSFTPCollectorSuite(t *testing.T) {
	suite.Run(t, new(SFTPCollectorSuite))
}

func (s *SFTPCollectorSuite) SetupTest() {
	log.SetLevel(log.DebugLevel)
	s.ctrl = gomock.NewController(s.T())
	s.sftpClient = mocks.NewMockSFTPClient(s.ctrl)
	s.collector = NewSFTPCollector(s.sftpClient)
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
		`Desc{fqName: "sftp_objects_available", `+
			`help: "Number of objects in the path", constLabels: {}, variableLabels: [path]}`,
		objectCount.String(),
	)

	objectSize := <-ch
	s.Equal(
		`Desc{fqName: "sftp_objects_total_size_bytes", `+
			`help: "Total size of all the objects in the path", constLabels: {}, variableLabels: [path]}`,
		objectSize.String(),
	)
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldWriteUpMetric() {
	viper.Set(viperkeys.SFTPPaths, []string{})
	s.sftpClient.EXPECT().Connect().Return(nil)
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
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/path0", 0755)
	_ = memFs.MkdirAll("/path1", 0755)
	path0Walker := fs.WalkFS("/path0", memKrFs{memFs: memFs})
	path1Walker := fs.WalkFS("/path1", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpClient.EXPECT().StatVFS("/path0").Return(&sftp.StatVFS{Frsize: 10, Blocks: 1000, Bfree: 100}, nil)
	s.sftpClient.EXPECT().StatVFS("/path1").Return(&sftp.StatVFS{Frsize: 5, Blocks: 1000, Bfree: 500}, nil)
	s.sftpClient.EXPECT().Walk("/path0").Return(path0Walker)
	s.sftpClient.EXPECT().Walk("/path1").Return(path1Walker)
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

	<-ch
	<-ch
	<-ch
	<-ch
	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldNotWriteFSMetricsOnError() {
	viper.Set(viperkeys.SFTPPaths, []string{"/path0"})
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/path0", 0755)
	path0Walker := fs.WalkFS("/path0", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpClient.EXPECT().StatVFS("/path0").Return(nil, fmt.Errorf("failed to get VFS stats"))
	s.sftpClient.EXPECT().Walk("/path0").Return(path0Walker)
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	m1 := <-ch
	s.NotContains(m1.Desc().String(), "filesystem_total_space_bytes")
	s.NotContains(m1.Desc().String(), "filesystem_free_space_bytes")
	m2 := <-ch
	s.NotContains(m2.Desc().String(), "filesystem_total_space_bytes")
	s.NotContains(m2.Desc().String(), "filesystem_free_space_bytes")
	m3 := <-ch
	s.NotContains(m3.Desc().String(), "filesystem_total_space_bytes")
	s.NotContains(m3.Desc().String(), "filesystem_free_space_bytes")
	close(ch)

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldWriteObjectMetrics() {
	viper.Set(viperkeys.SFTPPaths, []string{"/path0", "/path1"})
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/path0/1/a", 0755)
	_ = afero.WriteFile(memFs, "/path0/0.txt", []byte("0"), 0644)
	_ = afero.WriteFile(memFs, "/path0/1/1.txt", []byte("1"), 0644)
	_ = afero.WriteFile(memFs, "/path0/1/a/1a.txt", []byte("1a"), 0644)
	_ = memFs.MkdirAll("/path1/empty-dir", 0755)
	_ = afero.WriteFile(memFs, "/path1/1.txt", []byte("helloworld"), 0644)
	path0Walker := fs.WalkFS("/path0", memKrFs{memFs: memFs})
	path1Walker := fs.WalkFS("/path1", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpClient.EXPECT().StatVFS("/path0").Return(&sftp.StatVFS{}, nil)
	s.sftpClient.EXPECT().StatVFS("/path1").Return(&sftp.StatVFS{}, nil)
	s.sftpClient.EXPECT().Walk("/path0").Return(path0Walker)
	s.sftpClient.EXPECT().Walk("/path1").Return(path1Walker)
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	<-ch
	<-ch
	<-ch
	<-ch
	<-ch
	metric := &dto.Metric{}
	var desc *prometheus.Desc

	objectCount1 := <-ch
	desc = objectCount1.Desc()
	_ = objectCount1.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_available", help: "Number of objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(3.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path0", metric.GetLabel()[0].GetValue())

	objectSize1 := <-ch
	desc = objectSize1.Desc()
	_ = objectSize1.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_total_size_bytes", help: "Total size of all the objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(4.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path0", metric.GetLabel()[0].GetValue())

	objectCount2 := <-ch
	desc = objectCount2.Desc()
	_ = objectCount2.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_available", help: "Number of objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(1.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path1", metric.GetLabel()[0].GetValue())

	objectSize2 := <-ch
	desc = objectSize2.Desc()
	_ = objectSize2.Write(metric)
	s.Equal(`Desc{fqName: "sftp_objects_total_size_bytes", help: "Total size of all the objects in the path", `+
		`constLabels: {}, variableLabels: [path]}`, desc.String())
	s.Equal(10.0, metric.GetGauge().GetValue())
	s.Equal("path", metric.GetLabel()[0].GetName())
	s.Equal("/path1", metric.GetLabel()[0].GetValue())

	<-done
}

func (s *SFTPCollectorSuite) TestSFTPCollectorCollectShouldNotWriteObjectMetricsOnError() {
	viper.Set(viperkeys.SFTPPaths, []string{"/errorpath"})
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/errorpath", 0755)
	_ = afero.WriteFile(memFs, "/errorpath/file.txt", []byte("helloworld"), 0000)
	walker := fs.WalkFS("/errorpath", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Connect().Return(nil)
	s.sftpClient.EXPECT().StatVFS("/errorpath").Return(&sftp.StatVFS{}, nil)
	s.sftpClient.EXPECT().Walk("/errorpath").Return(walker)
	s.sftpClient.EXPECT().Close()
	ch := make(chan prometheus.Metric)
	done := make(chan bool)

	go func() {
		s.collector.Collect(ch)
		done <- true
	}()

	m1 := <-ch
	s.NotContains(m1.Desc().String(), "objects_available")
	s.NotContains(m1.Desc().String(), "objects_total_size_bytes")
	m2 := <-ch
	s.NotContains(m2.Desc().String(), "objects_available")
	s.NotContains(m2.Desc().String(), "objects_total_size_bytes")
	m3 := <-ch
	s.NotContains(m3.Desc().String(), "objects_available")
	s.NotContains(m3.Desc().String(), "objects_total_size_bytes")
	close(ch)

	<-done
}
