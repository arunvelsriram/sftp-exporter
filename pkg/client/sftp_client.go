package client

import (
	"fmt"

	"github.com/arunvelsriram/sftp-exporter/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/arunvelsriram/sftp-exporter/pkg/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type (
	SFTPClient interface {
		Close()
		FSStats() (model.FSStats, error)
		ObjectStats() (model.ObjectStats, error)
	}

	defaultSFTPClient struct {
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		config     config.Config
	}
)

func (d defaultSFTPClient) Close() {
	if err := d.sftpClient.Close(); err != nil {
		log.WithFields(log.Fields{
			"event": "closing SFTP connection"},
		).Error(err)
	}
	if err := d.sshClient.Close(); err != nil {
		log.WithFields(log.Fields{
			"event": "closing SSH connection"},
		).Error(err)
	}
}

func (d defaultSFTPClient) FSStats() (model.FSStats, error) {
	paths := d.config.GetSFTPPaths()
	fsStats := make([]model.FSStat, len(paths))
	for i, path := range paths {
		statVFS, err := d.sftpClient.StatVFS(path)
		if err != nil {
			return nil, err
		}
		fsStats[i] = model.FSStat{
			Path:       path,
			TotalSpace: float64(statVFS.TotalSpace()),
			FreeSpace:  float64(statVFS.FreeSpace()),
		}
	}
	return fsStats, nil
}

func (d defaultSFTPClient) ObjectStats() (model.ObjectStats, error) {
	paths := d.config.GetSFTPPaths()
	objectStats := make([]model.ObjectStat, len(paths))
	for i, path := range paths {
		walker := d.sftpClient.Walk(path)
		var size int64
		count := 0
		for walker.Step() {
			if err := walker.Err(); err != nil {
				log.WithFields(log.Fields{
					"event": "collecting object stats",
					"path":  path,
				}).Error(err)
				continue
			}

			if walker.Stat().IsDir() {
				continue
			}
			size += walker.Stat().Size()
			count++
		}

		objectStats[i] = model.ObjectStat{
			Path:        path,
			ObjectCount: float64(count),
			ObjectSize:  float64(size),
		}
	}

	return objectStats, nil
}

func NewSFTPClient(cfg config.Config) (SFTPClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.GetSFTPHost(), cfg.GetSFTPPort())
	auth, err := utils.SSHAuthMethods(cfg.GetSFTPPass(), cfg.GetSFTPKey(), cfg.GetSFTPKeyPassphrase())
	if err != nil {
		log.Error("unable to get SSH auth methods")
		return nil, err
	}
	clientConfig := &ssh.ClientConfig{
		User:            cfg.GetSFTPUser(),
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		if err := sshClient.Close(); err != nil {
			log.WithFields(log.Fields{
				"event": "closing SFTP connection"},
			).Error(err)
		}
		return nil, err
	}

	return defaultSFTPClient{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		config:     cfg,
	}, nil
}
