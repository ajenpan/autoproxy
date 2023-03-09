package main

import (
	"autoproxy/pool"
	"net/http"
)

func main() {
	pool := pool.NewPoolApp()
	// err := pool.Init()
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	mux := pool.ServeMux()

	// http.Handle("/pool/", mux)
	http.ListenAndServe(":8088", mux)
}
