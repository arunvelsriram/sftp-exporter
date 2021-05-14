package utils_test

import (
	"fmt"
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/utils"
	"github.com/stretchr/testify/assert"
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
