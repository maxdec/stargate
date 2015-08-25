package main

import (
	"net/http/httputil"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/transport"
)

// CtxBSON represents a Meta document in Mongo
type CtxBSON struct {
	Req  string //ReqBSON
	Resp string //RespBSON
	// time time.Time
	From string
}

// ReqBSON represents a Request in Mongo (nested in Meta)
type ReqBSON struct{}

// ResBSON represents a Response in Mongo (nested in Meta)
type ResBSON struct{}

// ProxyCtxToBSON transforms a ProxyCtx into a CtxBSON
func ProxyCtxToBSON(ctx *goproxy.ProxyCtx) (CtxBSON, error) {
	var ctxBSON CtxBSON

	reqBuf, err := httputil.DumpRequestOut(ctx.Req, false)
	if err != nil {
		return ctxBSON, err
	}
	ctxBSON.Req = string(reqBuf)

	respBuf, err := httputil.DumpResponse(ctx.Resp, false)
	if err != nil {
		return ctxBSON, err
	}
	ctxBSON.Resp = string(respBuf)

	ctxBSON.From = ""
	if ctx.UserData != nil {
		ctxBSON.From = ctx.UserData.(*transport.RoundTripDetails).TCPAddr.String()
	}

	return ctxBSON, nil
}
