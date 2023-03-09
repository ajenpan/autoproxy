package main

import (
	"autoproxy/log"
	"autoproxy/pool"
	"net/http"
)

func main() {
	pool := pool.NewPoolApp()
	err := pool.Init()
	if err != nil {
		log.Error(err)
		return
	}

	http.Handle("/pool/", pool)
	http.ListenAndServe(":8088", nil)
}
