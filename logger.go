package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"gopkg.in/mgo.v2"
)

var emptyResp = &http.Response{}
var emptyReq = &http.Request{}

// HTTPLogger saves Requests or Responses into the database
type HTTPLogger struct {
	ctxch     chan *goproxy.ProxyCtx
	errch     chan error
	dbsession *mgo.Session
}

// NewHTTPLogger creates a new HTTPLogger
func NewHTTPLogger(collection string) (*HTTPLogger, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	logger := &HTTPLogger{make(chan *goproxy.ProxyCtx), make(chan error), session}
	go logger.logAsync()

	return logger, nil
}

// LogAsync starts polling the ctx chanel
func (logger *HTTPLogger) logAsync() {
	for ctx := range logger.ctxch {
		logger.convertsAndSaves(ctx)
	}
	logger.dbsession.Close()
}

// LogCtx logs a ProxyCtx
func (logger *HTTPLogger) LogCtx(ctx *goproxy.ProxyCtx) {
	logger.ctxch <- ctx
}

func (logger *HTTPLogger) convertsAndSaves(ctx *goproxy.ProxyCtx) {
	ctxBSON, err := ProxyCtxToBSON(ctx)
	if err != nil {
		log.Println("Can't convert ctx into BSON:", err)
	}
	err = logger.dbsession.DB("stargate").C("reqres").Insert(ctxBSON)
	if err != nil {
		log.Println("Can't write ctx:", err, ctxBSON)
	}
}
