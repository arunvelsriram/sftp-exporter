package client

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
	"github.com/arunvelsriram/sftp-exporter/pkg/internal/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestSSHAuthMethods(t *testing.T) {
	tests := []struct {
		desc          string
		password      string
		key           string
		keyPassphrase string
		authMethods   []ssh.AuthMethod
		err           error
	}{
		{
			desc:          "should return error when key with invalid encoding is provided",
			password:      "",
			key:           "key-invalid-encoded",
			keyPassphrase: "",
			authMethods:   nil,
			err:           base64.CorruptInputError(3),
		},
		{
			desc:          "should return auth methods when password and key are given",
			password:      "password",
			key:           mocks.EncodedSSHKeyWithoutPassphrase(),
			keyPassphrase: "",
			authMethods:   []ssh.AuthMethod{ssh.PublicKeys(), ssh.Password("password")},
			err:           nil,
		},
		{
			desc:          "should get auth method when password is given",
			password:      "password",
			key:           "",
			keyPassphrase: "",
			authMethods:   []ssh.AuthMethod{ssh.Password("password")},
			err:           nil,
		},
		{
			desc:          "should return auth method when key is given",
			password:      "",
			key:           mocks.EncodedSSHKeyWithoutPassphrase(),
			keyPassphrase: "",
			authMethods:   []ssh.AuthMethod{ssh.PublicKeys()},
			err:           nil,
		},
		{
			desc:          "should return error when both password and key are empty",
			password:      "",
			key:           "",
			keyPassphrase: "",
			authMethods:   nil,
			err:           fmt.Errorf("failed to determine the SSH authentication methods to use"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			viper.Set(viperkeys.SFTPPassword, test.password)
			viper.Set(viperkeys.SFTPKey, test.key)
			viper.Set(viperkeys.SFTPKeyPassphrase, test.keyPassphrase)

			authMethods, err := sshAuthMethods()

			assert.Len(t, authMethods, len(test.authMethods))
			for i, expectedAuthMethod := range test.authMethods {
				expected := reflect.ValueOf(expectedAuthMethod).Type().Name()
				actual := reflect.ValueOf(authMethods[i]).Type().Name()
				assert.Equal(t, expected, actual)
			}
			assert.Equal(t, test.err, err)
		})
	}
}

func TestParsePrivateKey(t *testing.T) {
	tests := []struct {
		desc          string
		key           []byte
		keyPassphrase []byte
		err           error
	}{
		{
			desc:          "should parse key",
			key:           mocks.SSHKeyWithoutPassphrase(),
			keyPassphrase: []byte{},
			err:           nil,
		},
		{
			desc:          "should parse encrypted key",
			key:           mocks.SSHKeyWithPassphrase(),
			keyPassphrase: []byte(mocks.KeyPassphrase),
			err:           nil,
		},
		{
			desc:          "should return when invalid key is given",
			key:           []byte("invalid-key"),
			keyPassphrase: []byte(""),
			err:           fmt.Errorf("ssh: no key found"),
		},
		{
			desc:          "should return error when wrong passphrase is given",
			key:           mocks.SSHKeyWithPassphrase(),
			keyPassphrase: []byte("invalid-passphrase"),
			err:           fmt.Errorf("x509: decryption password incorrect"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_, err := parsePrivateKey(test.key, test.keyPassphrase)

			assert.Equal(t, test.err, err)
		})
	}
}
