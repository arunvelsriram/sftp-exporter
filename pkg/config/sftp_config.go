package config

import "fmt"

type SFTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

func (c SFTPConfig) toConfigMap() Configmap {
	return Configmap{
		"host": c.Host,
		"port": fmt.Sprintf("%d", c.Port),
		"user": c.User,
		"pass": c.Pass,
	}
}
