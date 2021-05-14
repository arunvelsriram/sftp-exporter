package client

import (
	"encoding/base64"
	"fmt"

	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func parsePrivateKey(key, keyPassphrase []byte) (parsedKey ssh.Signer, err error) {
	if len(keyPassphrase) > 0 {
		log.Debug("key has passphrase")
		parsedKey, err = ssh.ParsePrivateKeyWithPassphrase(key, keyPassphrase)
		if err != nil {
			log.WithField("when", "parsing encrypted ssh key").
				Error("failed to parse key with passphrase")
			return nil, err
		}
		return parsedKey, err
	}

	log.Debug("key has no passphrase")
	parsedKey, err = ssh.ParsePrivateKey(key)
	if err != nil {
		log.WithField("when", "parsing ssh key").Error("failed to parse key")
		return nil, err
	}
	return parsedKey, err
}

func sshAuthMethods() ([]ssh.AuthMethod, error) {
	password := viper.GetString(viperkeys.SFTPPassword)
	encodedKey := viper.GetString(viperkeys.SFTPKey)
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}
	keyPassphrase := []byte(viper.GetString(viperkeys.SFTPKeyPassphrase))

	if len(password) > 0 && len(key) > 0 {
		log.Debug("key and password are provided")
		parsedKey, err := parsePrivateKey(key, keyPassphrase)
		if err != nil {
			log.WithField("when", "determining SSH authentication methods").Error(err)
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
			ssh.Password(password),
		}, nil

	} else if len(password) > 0 {
		log.Debug("password is provided")
		return []ssh.AuthMethod{
			ssh.Password(password),
		}, nil
	} else if len(key) > 0 {
		log.Debug("key is provided")
		parsedKey, err := parsePrivateKey(key, keyPassphrase)
		if err != nil {
			log.WithField("when", "determining SSH authentication methods").Error(err)
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
		}, nil
	}

	log.Debug("both password and key are not provided")
	return nil, fmt.Errorf("failed to determine the SSH authentication methods to use")
}

func NewSSHClient() (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", viper.GetString(viperkeys.SFTPHost), viper.GetInt(viperkeys.SFTPPort))
	auth, err := sshAuthMethods()
	if err != nil {
		log.WithField("when", "creating a SSH client").Error(err)
		return nil, err
	}
	clientConfig := &ssh.ClientConfig{
		User:            viper.GetString(viperkeys.SFTPUser),
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", addr, clientConfig)
}
