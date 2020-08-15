package utils

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func SSHAuthMethods(pass string, key, keyPassphrase []byte) ([]ssh.AuthMethod, error) {
	if len(key) > 0 && IsNotEmpty(pass) {
		log.Debug("will be authenticating using key and password")
		parsedKey, err := parsePrivateKey(key, keyPassphrase)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
			ssh.Password(pass),
		}, nil
	} else if len(key) > 0 {
		log.Debug("will be authenticating using key")
		parsedKey, err := parsePrivateKey(key, keyPassphrase)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
		}, nil
	} else if IsNotEmpty(pass) {
		log.Debug("will be authenticating using password")
		return []ssh.AuthMethod{
			ssh.Password(pass),
		}, nil
	}
	return nil, fmt.Errorf("either one of password or key is required")
}

func parsePrivateKey(key, keyPassphrase []byte) (parsedKey ssh.Signer, err error) {
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
