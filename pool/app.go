package pool

import (
	"net/http"
	"net/http/pprof"
)

type PoolApp struct {
}

func (p *PoolApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do something

	pprof.Index(w, r)
}
