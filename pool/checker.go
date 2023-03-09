package pool

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CheckReport struct {
	Addr *ProxyAddr
	Cost time.Duration
	// Extra map[string]interface{}

	RespCode   int
	RespBody   []byte
	RespHeader http.Header
}

type Checker struct {
	Method    string
	TargetUrl string
	Client    *http.Client
}

func (c *Checker) HealthCheck(addr *ProxyAddr) (*CheckReport, error) {
	if addr == nil {
		return nil, fmt.Errorf("addr is nil")
	}

	if c.Method == "" {
		c.Method = "GET"
	}
	if c.TargetUrl == "" {
		c.TargetUrl = "https://httpbin.org/get"
	}
	if c.Client == nil {
		c.Client = &http.Client{}
	}

	c.Client.Transport = &http.Transport{
		Proxy:              http.ProxyURL(addr.URL()),
		ProxyConnectHeader: addr.Header.Clone(),
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	startAt := time.Now()
	req, err := http.NewRequest(c.Method, c.TargetUrl, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	report := &CheckReport{
		Addr:       addr,
		Cost:       time.Since(startAt),
		RespCode:   resp.StatusCode,
		RespBody:   respRaw,
		RespHeader: resp.Header,
	}
	return report, nil
}
