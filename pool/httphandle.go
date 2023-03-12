package pool

import (
	"encoding/json"
	"fmt"
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
	readReq := func(r *http.Request, out interface{}) error {
		if r.Method != "POST" {
			return fmt.Errorf("method not allowed")
		}

		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(raw, out)
		if err != nil {
			return err
		}
		return nil
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
	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		addrs := []string{}
		err := readReq(r, &addrs)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if len(addrs) > 0 {
			go func() {
				for _, addr := range addrs {
					h.App.DoCheck(addr)
				}
			}()
		}
		doresp(w, &RespType{Msg: "checking"}, nil)
	})

	mux.HandleFunc("/putin", func(w http.ResponseWriter, r *http.Request) {
		addrs := []string{}
		err := readReq(r, &addrs)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if len(addrs) > 0 {
			go h.App.Putin(addrs)
		}
		doresp(w, &RespType{Msg: "ok"}, nil)
	})
}
