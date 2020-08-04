package config

type SFTPConfig map[string]string

func (c SFTPConfig) Get(key string) (string, bool) {
	value, ok := c[key]
	return value, ok
}

func (c SFTPConfig) Set(key, value string) {
	c[key] = value
}
