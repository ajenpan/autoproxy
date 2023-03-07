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
		Addr : addr,
		Cost:       time.Since(startAt),
		RespCode:   resp.StatusCode,
		RespBody:   respRaw,
		RespHeader: resp.Header,
	}
	return report, nil
}

// func ListCheck(rows ) []*AddrModel { ) {
// 	raw, err := os.ReadFile("proxy.json")

// 	if err != nil {
// 		return
// 	}

// 	rows := make([]map[string]interface{}, 0)
// 	err = json.Unmarshal(raw, &rows)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	models := []*AddrModel{}

// 	type reporter struct {
// 		proxyAddr *ProxyAddr
// 		report    *CheckReport
// 		err       error
// 	}

// 	que := make(chan *reporter, len(rows))
// 	wg := &sync.WaitGroup{}

// 	for _, row := range rows {
// 		schema := row["schema"].(string)
// 		proxy := row["proxy"].(string)

// 		addrstr := strings.ToLower(schema) + "://" + proxy
// 		addr := NewProxyAddr(addrstr)
// 		if err != nil {
// 			continue
// 		}

// 		wg.Add(1)
// 		go func(addr *ProxyAddr) {
// 			defer wg.Done()

// 			checker := &Checker{
// 				TargetUrl: "https://www.baidu.com",
// 			}
// 			report, err := checker.HealthCheck(addr)
// 			que <- &reporter{
// 				proxyAddr: addr,
// 				report:    report,
// 				err:       err,
// 			}
// 		}(addr)
// 	}
// 	wg.Wait()

// 	close(que)

// 	for r := range que {
// 		if r.err != nil {
// 			fmt.Println(r.err, r.proxyAddr)
// 			continue
// 		}

// 		report := r.report
// 		condi := (report != nil && report.RespCode == 200)
// 		if !condi {
// 			fmt.Println("report code:", report.RespCode)
// 			continue
// 		}

// 		header, _ := json.Marshal(r.proxyAddr.Header)
// 		models = append(models, &AddrModel{
// 			Addr:   r.proxyAddr.Address(),
// 			Header: header,
// 			Health: 1,
// 			Speed:  r.report.Cost.Milliseconds(),
// 		})
// 	}

// 	store := &FileStorager{
// 		FileName: "proxies.json",
// 	}

// 	err = store.Save(models)

// 	if err != nil {
// 		t.Error(err)
// 	}
// }