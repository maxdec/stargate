package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func startProxy(addr string, verbose bool) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	logger, err := NewHTTPLogger("reqres")
	if err != nil {
		log.Fatal(err)
		return
	}

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		logger.LogCtx(ctx)
		return resp
	})

	log.Println("Starting proxy on address", addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("l", ":8080", "on which address should the proxy listen")
	flag.Parse()

	startProxy(*addr, *verbose)
}
