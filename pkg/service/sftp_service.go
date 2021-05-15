package service

import (
	"github.com/arunvelsriram/sftp-exporter/pkg/client"
	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/arunvelsriram/sftp-exporter/pkg/model"
	log "github.com/sirupsen/logrus"
)

type (
	SFTPService interface {
		ObjectStats() model.ObjectStats
	}

	sftpService struct {
		sftpClient client.SFTPClient
		config     config.Config
	}
)

func (s sftpService) ObjectStats() model.ObjectStats {
	paths := s.config.GetSFTPPaths()
	objectStats := make([]model.ObjectStat, 0)
	for _, path := range paths {
		walker := s.sftpClient.Walk(path)
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
		objectStat := model.ObjectStat{
			Path:        path,
			ObjectCount: float64(count),
			ObjectSize:  float64(size),
		}
		objectStats = append(objectStats, objectStat)
	}

	return objectStats
}

func NewSFTPService(cfg config.Config, s client.SFTPClient) SFTPService {
	return sftpService{
		sftpClient: s,
		config:     cfg,
	}
}
