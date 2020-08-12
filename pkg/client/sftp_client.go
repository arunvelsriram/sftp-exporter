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
		FSStat() (model.FSStats, error)
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

func (d defaultSFTPClient) FSStat() (model.FSStats, error) {
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
	return model.FSStats(fsStats), nil
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

	return model.ObjectStats(objectStats), nil
}

func ParsePrivateKey(key, keyPassphrase []byte) (parsedKey ssh.Signer, err error) {
	if len(keyPassphrase) > 0 {
		log.Debug("key has passphrase")
		parsedKey, err = ssh.ParsePrivateKeyWithPassphrase(key, keyPassphrase)
		if err != nil {
			log.Error("failed to parse key with passphrase")
			return nil, err
		}
		return parsedKey, err
	}

	log.Debug("key has no passphrase")
	parsedKey, err = ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error("failed to parse key")
		return nil, err
	}
	return parsedKey, err
}

func SSHAuthMethods(pass string, key, keyPassphrase []byte) ([]ssh.AuthMethod, error) {
	if len(key) > 0 && utils.IsNotEmpty(pass) {
		log.Debug("will be authenticating using key and password")
		parsedKey, err := ParsePrivateKey(key, keyPassphrase)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
			ssh.Password(pass),
		}, nil
	} else if len(key) > 0 {
		log.Debug("will be authenticating using key")
		parsedKey, err := ParsePrivateKey(key, keyPassphrase)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
		}, nil
	} else if utils.IsNotEmpty(pass) {
		log.Debug("will be authenticating using password")
		return []ssh.AuthMethod{
			ssh.Password(pass),
		}, nil
	}

	return nil, fmt.Errorf("either one of password or key is required")
}

func NewSFTPClient(cfg config.Config) (SFTPClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.GetSFTPHost(), cfg.GetSFTPPort())
	auth, err := SSHAuthMethods(cfg.GetSFTPPass(), cfg.GetSFTPKey(), cfg.GetSFTPKeyPassphrase())
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
