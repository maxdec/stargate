package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/transport"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var emptyResp = &http.Response{}
var emptyReq = &http.Request{}

// HTTPLogger saves Requests or Responses into the database
type HTTPLogger struct {
	metach    chan *Meta
	errch     chan error
	dbsession *mgo.Session
}

// NewHTTPLogger creates a new HTTPLogger
func NewHTTPLogger(collection string) (*HTTPLogger, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	logger := &HTTPLogger{make(chan *Meta), make(chan error), session}
	go logger.logAsync()

	return logger, nil
}

// LogAsync starts polling the meta chanel
func (logger *HTTPLogger) logAsync() {
	for m := range logger.metach {
		if err := m.Save(logger.dbsession.DB("stargate").C("reqres")); err != nil {
			log.Println("Can't write meta", err)
		}
	}
	logger.dbsession.Close()
}

// LogReq saves a Request
func (logger *HTTPLogger) LogReq(req *http.Request, ctx *goproxy.ProxyCtx) {
	if req == nil {
		req = emptyReq
	}
	logger.LogMeta(&Meta{
		id:   bson.NewObjectId(),
		req:  req,
		err:  ctx.Error,
		time: time.Now(),
		from: req.RemoteAddr,
	})
}

// LogRes saves a Response
func (logger *HTTPLogger) LogRes(res *http.Response, ctx *goproxy.ProxyCtx) {
	from := ""
	if ctx.UserData != nil {
		from = ctx.UserData.(*transport.RoundTripDetails).TCPAddr.String()
	}
	if res == nil {
		res = emptyResp
	}
	logger.LogMeta(&Meta{
		id:   bson.NewObjectId(),
		res:  res,
		err:  ctx.Error,
		time: time.Now(),
		from: from,
	})
}

// LogMeta enqueues the meta into the meta channel of the logger
func (logger *HTTPLogger) LogMeta(m *Meta) {
	logger.metach <- m
}
