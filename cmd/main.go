package main

import (
	"autoproxy/pool"
	"net/http"
)

func main() {
	app := pool.NewPoolApp()
	defer app.Stop()
	go app.Run()

	pool.NewHttpHandler(app, http.DefaultServeMux)
	http.ListenAndServe(":8088", nil)
}
