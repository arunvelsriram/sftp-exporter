package collector

import (
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/client"
	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/arunvelsriram/sftp-exporter/pkg/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/suite"
)

type SFTPCollectorSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	config         config.Config
	sftpClient     client.SFTPClient
	createClientFn CreateClientFn
	collector      prometheus.Collector
}

func TestSFTPCollectorSuite(t *testing.T) {
	suite.Run(t, new(SFTPCollectorSuite))
}

func (s *SFTPCollectorSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.config = mocks.NewMockConfig(s.ctrl)
	s.sftpClient = mocks.NewMockSFTPClient(s.ctrl)
	fn := func(c config.Config) (client.SFTPClient, error) {
		return s.sftpClient, nil
	}
	s.createClientFn = CreateClientFn(fn)
	s.collector = NewSFTPCollector(s.config, s.createClientFn)
}

func (s *SFTPCollectorSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SFTPCollectorSuite) TestSFTPCollectorDescribe() {
	ch := make(chan *prometheus.Desc)
	go s.collector.Describe(ch)

	up := <-ch
	s.Equal(
		`Desc{fqName: "sftp_up", help: "Tells if exporter is able to connect to SFTP", constLabels: {}, variableLabels: []}`,
		up.String(),
	)

	fsTotalSpace := <-ch
	s.Equal(
		`Desc{fqName: "sftp_filesystem_total_space_bytes", help: "Total space in the filesystem containing the path", constLabels: {}, variableLabels: [path]}`,
		fsTotalSpace.String(),
	)

	fsFreeSpace := <-ch
	s.Equal(
		`Desc{fqName: "sftp_filesystem_free_space_bytes", help: "Free space in the filesystem containing the path", constLabels: {}, variableLabels: [path]}`,
		fsFreeSpace.String(),
	)

	objectCount := <-ch
	s.Equal(
		`Desc{fqName: "sftp_objects_count_total", help: "Total number of objects in the path", constLabels: {}, variableLabels: [path]}`,
		objectCount.String(),
	)

	objectSize := <-ch
	s.Equal(
		`Desc{fqName: "sftp_objects_size_total_bytes", help: "Total size of all objects in the path", constLabels: {}, variableLabels: [path]}`,
		objectSize.String(),
	)
}
