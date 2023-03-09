package pool

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"
)

type PoolApp struct {
	addrslocker sync.RWMutex
	addrs       []string

	cronImp *cron.Cron

	frashLock sync.Mutex
	frashTask cron.EntryID
}

func NewPoolApp() *PoolApp {
	return &PoolApp{
		cronImp: cron.New(),
	}
}

func (app *PoolApp) ServeMux() *http.ServeMux {
	type RespType struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	doresp := func(w http.ResponseWriter, resp *RespType, err error) {
		if err != nil {
			resp.Code = -1
			resp.Msg = err.Error()
		}
		data, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/randone", func(w http.ResponseWriter, r *http.Request) {
		doresp(w, &RespType{
			Data: app.RandOnde(),
		}, nil)
	})

	mux.HandleFunc("/getall", func(w http.ResponseWriter, r *http.Request) {
		doresp(w, &RespType{
			Data: app.GetAll(),
		}, nil)
	})

	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		doresp(w, &RespType{
			Data: app.GetAll(),
		}, nil)
	})

	mux.HandleFunc("/frash", func(w http.ResponseWriter, r *http.Request) {
		go app.Flash()
		doresp(w, &RespType{Msg: "ok"}, nil)
	})

	mux.HandleFunc("/putin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		_, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//todo
	})

	return mux
}

func (app *PoolApp) Init() error {
	var err error
	app.frashTask, err = app.cronImp.AddFunc("0 0/10 * * * ?", func() {
		app.Flash()
	})
	app.Flash()
	return err
}

func (app *PoolApp) GetAll() []string {
	app.addrslocker.RLock()
	defer app.addrslocker.RUnlock()
	ret := app.addrs[:]
	return ret
}
func (app *PoolApp) AddrCount() int {
	app.addrslocker.RLock()
	defer app.addrslocker.RUnlock()
	return len(app.addrs)
}

func (app *PoolApp) RandOnde() string {
	app.addrslocker.RLock()
	defer app.addrslocker.RUnlock()

	cnt := len(app.addrs)
	if cnt == 0 {
		return ""
	}
	i := rand.Int31n(int32(cnt))
	return app.addrs[i]
}

func (app *PoolApp) Flash() {
	ok := app.frashLock.TryLock()
	if !ok {
		return
	}
	defer app.frashLock.Unlock()

	// oldaddrs := app.GetAll()
	// defer app.frashLock.Unlock()

	// cache := make(map[string]struct{})
	// for _, addrstr := range oldaddrs {
	// 	cache[addrstr] = struct{}{}
	// }

	// crawlResult := crawl.RunAllCrawlers()

	// for _, item := range crawlResult {
	// 	for _, addrstr := range item.Addrs {
	// 		cache[addrstr] = struct{}{}
	// 	}
	// }

	// addrs := make([]*ProxyAddr, 0, len(cache))
	// for k := range cache {
	// 	addr, err := NewProxyAddrFromStr(k)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	addrs = append(addrs, addr)
	// }

	// addrs = ListHealth2(addrs)

	// newaddrs := make([]string, 0, len(addrs))
	// for _, addr := range addrs {
	// 	newaddrs = append(newaddrs, addr.Address())
	// }

	// app.addrslocker.Lock()
	// defer app.addrslocker.Unlock()
	// app.addrs = newaddrs
}
