package pool

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"

	"autoproxy/pool/crawl"
)

type PoolApp struct {
	addrslocker sync.RWMutex
	addrs       []string
	cronImp     *cron.Cron
}

func NewPoolApp() *PoolApp {
	return &PoolApp{
		cronImp: cron.New(),
	}
}

func (p *PoolApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do something
}

func (app *PoolApp) Init() error {
	app.cronImp.AddFunc("0 0/10 * * * ?", func() {
		app.Flash()
	})
	app.Flash()
	return nil
}

func (app *PoolApp) GetAll() []string {
	app.addrslocker.Lock()
	defer app.addrslocker.Unlock()
	ret := app.addrs[:]
	return ret
}

func (app *PoolApp) Flash() {
	oldaddrs := app.GetAll()
	cache := make(map[string]struct{})
	for _, addrstr := range oldaddrs {
		cache[addrstr] = struct{}{}
	}

	crawlResult := crawl.RunAllCrawlers()

	for _, item := range crawlResult {
		for _, addrstr := range item.Addrs {
			cache[addrstr] = struct{}{}
		}
	}

	addrs := make([]*ProxyAddr, 0, len(cache))
	for k := range cache {
		addr, err := NewProxyAddrFromStr(k)
		if err != nil {
			fmt.Println(err)
			continue
		}
		addrs = append(addrs, addr)
	}

	addrs = ListHealth2(addrs)

	newaddrs := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		newaddrs = append(newaddrs, addr.Address())
	}

	app.addrslocker.Lock()
	defer app.addrslocker.Unlock()
	app.addrs = newaddrs
}
