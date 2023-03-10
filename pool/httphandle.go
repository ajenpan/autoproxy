package pool

import (
	"encoding/json"
	"io"
	"net/http"
)

type HttpHandle struct {
	App *PoolApp
}

func NewHttpHandler(app *PoolApp, mux *http.ServeMux) *HttpHandle {
	ret := &HttpHandle{
		App: app,
	}
	if mux == nil {
		mux = http.DefaultServeMux
	}
	ret.init(mux)
	return ret
}

func (h *HttpHandle) init(mux *http.ServeMux) {
	type RespType struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	doresp := func(w http.ResponseWriter, resp *RespType, err error) {
		if resp == nil {
			resp = &RespType{}
		}
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}

	mux.HandleFunc("/randone", func(w http.ResponseWriter, r *http.Request) {
		one, err := h.App.RandOnde()
		doresp(w, &RespType{
			Data: one,
		}, err)
	})

	mux.HandleFunc("/getall", func(w http.ResponseWriter, r *http.Request) {
		doresp(w, &RespType{
			Data: h.App.GetAll(),
		}, nil)
	})

	// mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
	// 	doresp(w, &RespType{}, nil)
	// })

	mux.HandleFunc("/frash", func(w http.ResponseWriter, r *http.Request) {
		go h.App.Frash()
		doresp(w, &RespType{Msg: "ok"}, nil)
	})

	mux.HandleFunc("/size", func(w http.ResponseWriter, r *http.Request) {
		doresp(w, &RespType{Data: h.App.Size()}, nil)
	})

	mux.HandleFunc("/putin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		addrs := []string{}

		err = json.Unmarshal(raw, &addrs)
		if err != nil {
			doresp(w, &RespType{}, err)
			return
		}

		go h.App.Putin(addrs)
		doresp(w, &RespType{Msg: "ok"}, nil)
	})
}
