package utils_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/internal/mocks"
	"github.com/arunvelsriram/sftp-exporter/pkg/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected bool
	}{
		{
			name:     "empty string",
			in:       "",
			expected: true,
		},
		{
			name:     "non-empty string",
			in:       "some value",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := utils.IsEmpty(test.in)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected bool
	}{
		{
			name:     "empty string",
			in:       "",
			expected: false,
		},
		{
			name:     "non-empty string",
			in:       "some value",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := utils.IsNotEmpty(test.in)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestPanicIfErrShouldPanicForErr(t *testing.T) {
	err := fmt.Errorf("some error")

	assert.PanicsWithValue(t, "some error", func() { utils.PanicIfErr(err) })
}

func TestPanicIfErrShouldNotPanicWhenErrIsNil(t *testing.T) {
	assert.NotPanics(t, func() { utils.PanicIfErr(nil) })
}

func TestSSHAuthMethods(t *testing.T) {
	tests := []struct {
		desc          string
		pass          string
		key           []byte
		keyPassphrase []byte
		authMethods   []ssh.AuthMethod
		err           error
	}{
		{
			desc:          "should return error when both pass and key are empty",
			pass:          "",
			key:           []byte{},
			keyPassphrase: []byte{},
			authMethods:   nil,
			err:           fmt.Errorf("either one of password or key is required"),
		},
		{
			desc:          "should return pass and key auth methods when pass and key are provided",
			pass:          "password",
			key:           mocks.SSHKeyWithoutPassphrase(),
			keyPassphrase: []byte{},
			authMethods:   []ssh.AuthMethod{ssh.PublicKeys(), ssh.Password("password")},
			err:           nil,
		},
		{
			desc:          "should return pass and key auth methods when pass and encrypted key are provided",
			pass:          "password",
			key:           mocks.SSHKeyWithPassphrase(),
			keyPassphrase: mocks.SSHKeyPassphrase(),
			authMethods:   []ssh.AuthMethod{ssh.PublicKeys(), ssh.Password("password")},
			err:           nil,
		},
		{
			desc:          "should return error when pass and invalid key are provided",
			pass:          "password",
			key:           []byte("invalid-key"),
			keyPassphrase: []byte{},
			authMethods:   nil,
			err:           fmt.Errorf("ssh: no key found"),
		},
		{
			desc:          "should return error when pass and wrong key passphrase are provided",
			pass:          "password",
			key:           mocks.SSHKeyWithPassphrase(),
			keyPassphrase: []byte("wrong-passphrase"),
			authMethods:   nil,
			err:           fmt.Errorf("x509: decryption password incorrect"),
		},
		{
			desc:          "should return key auth method when only key is provided",
			pass:          "",
			key:           mocks.SSHKeyWithoutPassphrase(),
			keyPassphrase: []byte{},
			authMethods:   []ssh.AuthMethod{ssh.PublicKeys()},
			err:           nil,
		},
		{
			desc:          "should return key auth method when only encrypted key is provided",
			pass:          "",
			key:           mocks.SSHKeyWithPassphrase(),
			keyPassphrase: mocks.SSHKeyPassphrase(),
			authMethods:   []ssh.AuthMethod{ssh.PublicKeys()},
			err:           nil,
		},
		{
			desc:          "should return error when key is invalid",
			pass:          "",
			key:           []byte("invalid-key"),
			keyPassphrase: []byte{},
			authMethods:   nil,
			err:           fmt.Errorf("ssh: no key found"),
		},
		{
			desc:          "should return error when key passphrase is wrong",
			pass:          "",
			key:           mocks.SSHKeyWithPassphrase(),
			keyPassphrase: []byte("wrong-passphrase"),
			authMethods:   nil,
			err:           fmt.Errorf("x509: decryption password incorrect"),
		},
		{
			desc:          "should return pass auth method when only pass is provided",
			pass:          "password",
			key:           []byte{},
			keyPassphrase: []byte{},
			authMethods:   []ssh.AuthMethod{ssh.Password("password")},
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			authMethods, err := utils.SSHAuthMethods(test.pass, test.key, test.keyPassphrase)

			assert.Len(t, authMethods, len(test.authMethods))
			for i, expectedAuthMethod := range test.authMethods {
				expected := runtime.FuncForPC(reflect.ValueOf(expectedAuthMethod).Pointer()).Name()
				actual := runtime.FuncForPC(reflect.ValueOf(authMethods[i]).Pointer()).Name()
				assert.Equal(t, expected, actual)
			}
			assert.Equal(t, test.err, err)
		})
	}
}
