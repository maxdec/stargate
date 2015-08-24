package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/transport"
)

func startProxy(addr string, verbose bool) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	logger, err := NewHTTPLogger("reqres")

	tr := transport.Transport{Proxy: transport.ProxyFromEnvironment}
	// For every incoming request, override the RoundTripper to extract
	// connection information. Store it is session context log it after
	// handling the response.
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {
			ctx.UserData, resp, err = tr.DetailedRoundTrip(req)
			return
		})
		// logger.LogReq(req, ctx)
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		// logger.LogResp(resp, ctx)
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
