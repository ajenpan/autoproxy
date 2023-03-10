package pool

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"autoproxy/log"
	"autoproxy/pool/crawl"
)

type PoolApp struct {
	cronImp *cron.Cron

	frashLock sync.Mutex
	frashTask cron.EntryID

	Source   Source
	Checker  Check
	Storager Storage
	log      log.Logger

	checkInterval time.Duration
}

func NewPoolApp() *PoolApp {
	app := &PoolApp{
		cronImp:       cron.New(),
		checkInterval: 10 * time.Minute,
		Source: NewSrouceGroup(
			crawl.NewCrawler(),
		),
		Checker:  NewRequestChecker("https://httpbin.org/get"),
		Storager: NewFileStorager("proxies.json"),
		log:      log.NewNoop(),
	}
	app.frashTask, _ = app.cronImp.AddFunc("*/10 * * * *", func() {
		app.Frash()
	})
	return app
}

func (app *PoolApp) Run() error {
	var err error
	go func() {
		for v := range app.Source.Reader() {
			app.Putin([]string{v})
		}
	}()
	go app.Frash()

	app.cronImp.Run()
	return err
}

func (app *PoolApp) Stop() {
	// better cannot stop during frash

	app.frashLock.Lock()
	defer app.frashLock.Unlock()

	app.cronImp.Stop()
	app.Source.Close()
	app.Storager.Flash()
}

func (app *PoolApp) GetAll() []*ProxyAddr {
	result, err := app.Storager.GetAll()
	if err != nil {
		return []*ProxyAddr{}
	}
	return result
}

func (app *PoolApp) Size() int {
	return app.Storager.Count()
}

func (app *PoolApp) RandOnde() (*ProxyAddr, error) {
	return app.Storager.RandOne()
}

func (app *PoolApp) Putin(addrstrs []string) {
	addrs := make([]*ProxyAddr, 0, len(addrstrs))
	for _, addrstr := range addrstrs {
		addr, _ := app.Storager.Get(addrstr)
		if addr != nil {
			continue
		}
		var err error
		addr, err = NewProxyAddrWithStr(addrstr)
		if err != nil {
			continue
		}
		addrs = append(addrs, addr)
	}

	for _, addr := range addrs {
		go func(addr *ProxyAddr) {
			fmt.Println("check", addr.Addr)
			app.DoHealthCheck(addr)
		}(addr)
	}
}

func (app *PoolApp) DoHealthCheck(addr *ProxyAddr) {
	app.Checker.HealthCheck(addr)
	if addr.Health >= -1 {
		app.AvailableAddr(addr)
	} else {
		app.UnavailableAddr(addr)
	}
}

func (app *PoolApp) AvailableAddr(addr *ProxyAddr) {
	app.Storager.Save(addr)
}

func (app *PoolApp) UnavailableAddr(addr *ProxyAddr) {
	app.Storager.Delete(addr.Addr)
}

func (app *PoolApp) Frash() {
	ok := app.frashLock.TryLock()
	if !ok {
		return
	}
	defer app.frashLock.Unlock()

	addrs := app.GetAll()
	needCheck := make([]*ProxyAddr, 0, len(addrs))

	for _, addr := range addrs {
		if time.Since(addr.CheckAt) > app.checkInterval {
			needCheck = append(needCheck, addr)
		}
	}
	wg := sync.WaitGroup{}
	for _, addr := range needCheck {
		wg.Add(1)
		go func(addr *ProxyAddr) {
			defer wg.Done()
			app.DoHealthCheck(addr)
		}(addr)
	}
	wg.Wait()
	app.Storager.Flash()
}
