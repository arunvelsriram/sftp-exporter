package utils

import (
	"strings"
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
