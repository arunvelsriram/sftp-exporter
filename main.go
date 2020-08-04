package main

import (
	"context"
	"fmt"

	"github.com/rclone/rclone/backend/sftp"
	"github.com/rclone/rclone/fs/config/obscure"
	"github.com/rclone/rclone/fs/operations"
)

type SFTPConfig map[string]string

func (c SFTPConfig) Get(key string) (string, bool) {
	value, ok := c[key]
	return value, ok
}

func (c SFTPConfig) Set(key, value string) {
	c[key] = value
}

func main() {
	pass, err := obscure.Obscure("pass")
	if err != nil {
		panic(err)
	}
	sftpConfig := SFTPConfig{"host": "localhost", "port": "22", "user": "foo", "pass": pass}
	fs, err := sftp.NewFs("my-sftp", "", sftpConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("got fs")

	size, count, err := operations.Count(context.Background(), fs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Size: %d\n", size)
	fmt.Printf("Count: %d\n", count)
}
