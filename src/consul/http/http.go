package consulhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ConsulKey = "/v1/kv/"
)

type API struct {
	Host       string
	Token      string
	Datacenter string
	Prefix     string
	Client     *http.Client
}

func New(host string, token string, dc string, prefix string) *API {
	var a API

	if host == "" {
		host = "http://localhost:8500"
	}

	a.Host = host
	a.Token = token
	a.Datacenter = dc
	a.Prefix = prefix

	a.Client = &http.Client{
		Timeout: 3 * time.Second,
	}

	// Test connection
	req, reqErr := http.NewRequest(http.MethodGet, a.URLKeyList("", ""), nil)
	if reqErr != nil {
		log.Fatalln("Couldn't create HTTPNewRequest.")
	}

	doResp, doErr := a.Client.Do(req)
	if doErr != nil {
		log.Fatalln("Couldn't create request to Consul API.")
	}

	defer doResp.Body.Close()

	if doResp.StatusCode != http.StatusOK {
		log.Fatalln("Consul API returning non-200 Status. Are you sure there's KV?")
	}

	return &a
}

func (a *API) URLKeyCreateUpdate(key string, token string) string {
	if a.Prefix != "" {
		key = fmt.Sprintf("%s%s", a.Prefix, key)
	}

	if token == "" {
		token = a.Token
	}

	return fmt.Sprintf("%s%s%s?token=%s&dc=%s", a.Host, ConsulKey, key, url.QueryEscape(token), url.QueryEscape(a.Datacenter))
}

func (a *API) URLKeyDelete(key string, token string) string {
	if a.Prefix != "" {
		key = fmt.Sprintf("%s%s", a.Prefix, key)
	}

	if token == "" {
		token = a.Token
	}

	recurse := ""
	if len(key) > 0 && key[len(key)-1] == '/' {
		recurse = "&recurse=true"
	}

	return fmt.Sprintf("%s%s%s?token=%s&dc=%s%s", a.Host, ConsulKey, key, url.QueryEscape(token), url.QueryEscape(a.Datacenter), recurse)
}

func (a *API) URLKeyGet(key string, token string) string {
	if a.Prefix != "" {
		key = fmt.Sprintf("%s%s", a.Prefix, key)
	}

	if token == "" {
		token = a.Token
	}

	return fmt.Sprintf("%s%s%s?token=%s&dc=%s&raw=true", a.Host, ConsulKey, key, url.QueryEscape(token), url.QueryEscape(a.Datacenter))
}

func (a *API) URLKeyList(key string, token string) string {
	if a.Prefix != "" {
		key = fmt.Sprintf("%s%s", a.Prefix, key)
	}

	if token == "" {
		token = a.Token
	}

	if a.Datacenter == "" {
		return fmt.Sprintf("%s%s%s?token=%s&keys=true&separator=/", a.Host, ConsulKey, key, a.Token)
	} else {
		return fmt.Sprintf("%s%s%s?token=%s&dc=%s&keys=true&separator=/", a.Host, ConsulKey, key, url.QueryEscape(token), url.QueryEscape(a.Datacenter))
	}
}

func (a *API) CreateUpdateKey(key string, val string, token string) bool {
	if key == "" {
		return false
	}

	req, reqErr := http.NewRequest(http.MethodPut, a.URLKeyCreateUpdate(key, token), strings.NewReader(val))
	if reqErr == nil {
		doResp, doErr := a.Client.Do(req)
		if doErr == nil {
			defer doResp.Body.Close()

			if doResp.StatusCode == http.StatusOK {
				bodyBytes, errReadBody := ioutil.ReadAll(doResp.Body)
				if errReadBody == nil {
					if string(bodyBytes) == "true" {
						return true
					}
				}
			}
		}
	}

	return false
}

func (a *API) DeleteKey(key string, token string) bool {
	if key == "" {
		return false
	}

	req, reqErr := http.NewRequest(http.MethodDelete, a.URLKeyDelete(key, token), nil)
	if reqErr == nil {
		doResp, doErr := a.Client.Do(req)
		if doErr == nil {
			defer doResp.Body.Close()

			if doResp.StatusCode == http.StatusOK {
				bodyBytes, errReadBody := ioutil.ReadAll(doResp.Body)
				if errReadBody == nil {
					if string(bodyBytes) == "true" {
						return true
					}
				}
			}
		}
	}

	return false
}

func (a *API) GetKeyValue(key string, token string) string {
	if key == "" {
		return ""
	}

	req, reqErr := http.NewRequest(http.MethodGet, a.URLKeyGet(key, token), nil)
	if reqErr == nil {
		doResp, doErr := a.Client.Do(req)
		if doErr == nil {
			defer doResp.Body.Close()

			if doResp.StatusCode == http.StatusOK {
				bodyBytes, errReadBody := ioutil.ReadAll(doResp.Body)
				if errReadBody == nil {
					return string(bodyBytes)
				}
			}
		}
	}

	return ""
}

func (a *API) GetListKey(key string, token string) []map[string]string {
	req, reqErr := http.NewRequest(http.MethodGet, a.URLKeyList(key, token), nil)

	if reqErr == nil {
		doResp, doErr := a.Client.Do(req)
		if doErr == nil {
			defer doResp.Body.Close()

			if doResp.StatusCode == http.StatusOK {
				bodyBytes, errReadBody := ioutil.ReadAll(doResp.Body)
				if errReadBody == nil {
					var keyList []string
					errUnmarshal := json.Unmarshal(bodyBytes, &keyList)

					if errUnmarshal == nil {
						var listFolder []map[string]string
						var listFile []map[string]string

						for _, keyVal := range keyList {
							if keyVal == key {
								continue
							}

							keyInterface := make(map[string]string)
							keyInterface["key"] = keyVal
							if a.Prefix != "" {
								keyInterface["key"] = keyInterface["key"][len(a.Prefix):]
							}

							if len(keyVal) > 0 && keyVal[len(keyVal)-1] == '/' {
								keyInterface["type"] = "folder"
								listFolder = append(listFolder, keyInterface)
							} else {
								keyInterface["type"] = "key"
								listFile = append(listFile, keyInterface)
							}
						}

						return append(listFolder, listFile...)
					}
				}
			}
		}
	}

	return nil
}
