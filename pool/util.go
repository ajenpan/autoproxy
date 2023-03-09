package pool

import (
	"encoding/json"
	"fmt"
	"sync"

	"autoproxy/pool/crawl"
)

func ListHealth2(addrs []*ProxyAddr) []*ProxyAddr {
	targetUrl := "https://httpbin.org/get"

	type reporter struct {
		report *CheckReport
		err    error
		addr   *ProxyAddr
	}

	que := make(chan *reporter, len(addrs))
	wg := &sync.WaitGroup{}
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr *ProxyAddr) {
			defer wg.Done()

			checker := &Checker{
				TargetUrl: targetUrl,
			}
			report, err := checker.HealthCheck(addr)
			que <- &reporter{
				addr:   addr,
				report: report,
				err:    err,
			}
		}(addr)
	}
	wg.Wait()
	close(que)

	ret := make([]*ProxyAddr, 0, len(addrs))

	for r := range que {
		if r.err != nil {
			fmt.Println(r.addr.Address(), r.err)
			continue
		}
		if r.report == nil {
			continue
		}
		report := r.report
		condi := (report != nil && report.RespCode == 200)
		if !condi {
			continue
		}
		ret = append(ret, r.addr)
	}
	return ret
}

func ListHealth(rows []*AddrModel) []*AddrModel {
	// if targetUrl == "" {
	// 	targetUrl = "https://www.baidu.com"
	// }
	// targetUrl := "https://www.baidu.com"
	targetUrl := "https://httpbin.org/get"

	type reporter struct {
		report *CheckReport
		err    error
		mod    *AddrModel
	}

	que := make(chan *reporter, len(rows))
	wg := &sync.WaitGroup{}
	for _, row := range rows {
		wg.Add(1)
		go func(row *AddrModel) {
			defer wg.Done()

			addr, err := NewProxyAddrFromStr(row.Addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(row.Header) > 2 {
				json.Unmarshal(row.Header, &addr.Header)
			}

			checker := &Checker{
				TargetUrl: targetUrl,
			}
			report, err := checker.HealthCheck(addr)
			que <- &reporter{
				mod:    row,
				report: report,
				err:    err,
			}
		}(row)
	}
	wg.Wait()
	close(que)

	ret := make([]*AddrModel, 0, len(rows))

	for r := range que {
		if r.err != nil {
			fmt.Println(r.report.Addr, r.err)
			continue
		}
		report := r.report
		condi := (report != nil && report.RespCode == 200)
		if !condi {
			fmt.Println("report code:", report.RespCode)
			continue
		}
		r.mod.Health = 1
		r.mod.Speed = report.Cost
		ret = append(ret, r.mod)
	}
	return ret
}

func UpdateSourceWithCrawler() {
	var err error
	crawlResult := crawl.RunAllCrawlers()
	cache := make(map[string]struct{})

	rows := []*AddrModel{}

	for _, item := range crawlResult {
		if item.Err != nil {
			fmt.Printf("crawl: %s error:%v \n", item.CrawlerName, item.Err)
			continue
		}
		for _, addrstr := range item.Addrs {
			if _, ok := cache[addrstr]; ok {
				continue
			}
			cache[addrstr] = struct{}{}
			rows = append(rows, &AddrModel{
				Addr: addrstr,
			})
		}
	}

	rows = ListHealth(rows)

	// store
	store := &FileStorager{
		FileName: "proxies.json",
	}

	oldrows, _ := store.Read()
	if len(oldrows) > 0 {
		rows = append(rows, oldrows...)
	}

	err = store.Save(rows)
	if err != nil {
		fmt.Println(err)
	}
}
