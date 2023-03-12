package pool

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Check interface {
	HealthCheck(addr *ProxyAddr) *CheckReport
}

type CheckReport struct {
	Addr  *ProxyAddr
	Extra interface{}
	Err   error
}

func NewRequestChecker(url string) Check {
	return &RequestChecker{
		Method:    "GET",
		TargetUrl: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type RequestChecker struct {
	Method    string
	TargetUrl string
	Client    *http.Client
	MaxHealth int8
}

type RequestCheckerExtra struct {
	RespCode   int
	RespBody   []byte
	RespHeader http.Header
}

func (c *RequestChecker) HealthCheck(addr *ProxyAddr) *CheckReport {
	report := &CheckReport{
		Addr: addr,
	}

	if addr == nil {
		report.Err = fmt.Errorf("addr is nil")
		return report
	}

	c.Client.Transport = &http.Transport{
		//ProxyConnectHeader: addr.Header,
		DisableKeepAlives:     true,
		Proxy:                 http.ProxyURL(addr.URL()),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   2 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	startAt := time.Now()
	req, err := http.NewRequest(c.Method, c.TargetUrl, nil)

	if err != nil {
		addr.Health = -1
		report.Err = err
		return report
	}

	resp, err := c.Client.Do(req)

	if err == nil && resp != nil && resp.StatusCode == 200 {
		addr.Health = 1

		extra := &RequestCheckerExtra{
			RespCode:   resp.StatusCode,
			RespHeader: resp.Header,
		}
		defer resp.Body.Close()
		extra.RespBody, _ = io.ReadAll(resp.Body)
		report.Extra = extra

	} else {
		addr.Health = -1
		report.Err = err
	}

	addr.Speed = time.Since(startAt).String()
	addr.CheckAt = time.Now()
	return report
}
