package service_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/internal/mocks"
	"github.com/arunvelsriram/sftp-exporter/pkg/model"
	"github.com/arunvelsriram/sftp-exporter/pkg/service"
	"github.com/golang/mock/gomock"
	"github.com/kr/fs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
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

type SFTPServiceSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	config     *mocks.MockConfig
	sftpClient *mocks.MockSFTPClient
	service    service.SFTPService
}

func TestSFTPServiceSuite(t *testing.T) {
	suite.Run(t, new(SFTPServiceSuite))
}

func (s *SFTPServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.config = mocks.NewMockConfig(s.ctrl)
	s.sftpClient = mocks.NewMockSFTPClient(s.ctrl)
	s.service = service.NewSFTPService(s.config, s.sftpClient)
}

func (s *SFTPServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SFTPServiceSuite) TestSFTPServiceObjectStats() {
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/path0/1", 0755)
	_ = memFs.MkdirAll("/path0/1/a", 0755)
	_ = memFs.MkdirAll("/path0/2", 0755)
	_ = afero.WriteFile(memFs, "/path0/0.txt", []byte("0"), 0644)
	_ = afero.WriteFile(memFs, "/path0/1/1.txt", []byte("1"), 0644)
	_ = afero.WriteFile(memFs, "/path0/1/a/1a.txt", []byte("1a"), 0644)
	_ = afero.WriteFile(memFs, "/path0/2/2.txt", []byte("2"), 0644)
	_ = memFs.MkdirAll("/path1", 0755)
	_ = afero.WriteFile(memFs, "/path1/1.txt", []byte("helloworld"), 0644)

	s.config.EXPECT().GetSFTPPaths().Return([]string{"/path0", "/path1"})
	path0Walker := fs.WalkFS("/path0", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Walk("/path0").Return(path0Walker)
	path1Walker := fs.WalkFS("/path1", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Walk("/path1").Return(path1Walker)

	objectStats := s.service.ObjectStats()

	expected := model.ObjectStats([]model.ObjectStat{
		{
			Path:        "/path0",
			ObjectCount: 4,
			ObjectSize:  5,
		},
		{
			Path:        "/path1",
			ObjectCount: 1,
			ObjectSize:  10,
		},
	})
	s.Equal(expected, objectStats)
}

func (s *SFTPServiceSuite) TestSFTPServiceObjectStatsShouldSkipDirs() {
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/path/dir1", 0755)
	_ = memFs.MkdirAll("/path/dir2", 0755)

	s.config.EXPECT().GetSFTPPaths().Return([]string{"/path"})
	walker := fs.WalkFS("/path", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Walk("/path").Return(walker)

	objectStats := s.service.ObjectStats()

	expected := model.ObjectStats([]model.ObjectStat{
		{
			Path:        "/path",
			ObjectCount: 0,
			ObjectSize:  0,
		},
	})
	s.Equal(expected, objectStats)
}

func (s *SFTPServiceSuite) TestSFTPServiceObjectStatsShouldSkipAndContinueInCaseOfError() {
	memFs := afero.NewMemMapFs()
	_ = memFs.MkdirAll("/errorpath", 0755)
	_ = afero.WriteFile(memFs, "/errorpath/file.txt", []byte("helloworld"), 0000)
	s.config.EXPECT().GetSFTPPaths().Return([]string{"/errorpath"})
	walker := fs.WalkFS("/errorpath", memKrFs{memFs: memFs})
	s.sftpClient.EXPECT().Walk("/errorpath").Return(walker)

	objectStats := s.service.ObjectStats()

	expected := model.ObjectStats([]model.ObjectStat{
		{
			Path:        "/errorpath",
			ObjectCount: 0,
			ObjectSize:  0,
		},
	})
	s.Equal(expected, objectStats)
}
