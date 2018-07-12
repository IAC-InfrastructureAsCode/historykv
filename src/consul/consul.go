package consul

import (
	historyConsulHTTP "github.com/tokopedia/historykv/src/consul/http"
)

type Consul interface {
	CreateUpdateKey(key string, val string, token string) bool
	DeleteKey(key string, token string) bool
	GetKeyValue(key string, token string) string
	GetListKey(key string, token string) []map[string]string
}

func New(host string, token string, dc string, prefix string) Consul {
	var c Consul

	if host == "" {
		host = "http://localhost:8500"
	}

	c = historyConsulHTTP.New(host, token, dc, prefix)

	return c
}
