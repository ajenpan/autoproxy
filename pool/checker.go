package pool

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type CheckReport struct {
	Cost  time.Duration
	Extra map[string]interface{}
}

type Checker struct {
	Method    string
	TargetUrl string
}

func (c *Checker) HealthCheck(addr *ProxyAddr) (*CheckReport, error) {
	if c.Method == "" {
		c.Method = "GET"
	}
	if c.TargetUrl == "" {
		c.TargetUrl = "https://httpbin.org/get"
	}
	addurl := addr.Url()

	var client = &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(addurl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	beginAt := time.Now()
	req, err := http.NewRequest(c.Method, c.TargetUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	report := &CheckReport{
		Cost:  time.Since(beginAt),
		Extra: make(map[string]interface{}),
	}

	json.Unmarshal(respRaw, &report.Extra)
	return report, nil
}
