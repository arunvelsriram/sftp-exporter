package config

type Configmap map[string]string

func (c Configmap) Get(key string) (string, bool) {
	value, ok := c[key]
	return value, ok
}

func (c Configmap) Set(key, value string) {
	c[key] = value
}
