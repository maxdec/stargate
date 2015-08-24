package main

import (
	"net/http"
	"net/http/httputil"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Enum is just a disguised int
type Enum int

const (
	// REQUEST represents a Request
	REQUEST Enum = iota
	// RESPONSE represents a Response
	RESPONSE
	// UNKNOWN ??? not a REQUEST or a RESPONSE => something went wrong
	UNKNOWN
)

// Meta contains the data for one half of a roundtrip (request or response)
type Meta struct {
	id   bson.ObjectId
	req  *http.Request
	res  *http.Response
	err  error
	time time.Time
	body string
	from string
}

// MetaBSON represents a Meta document in Mongo
type MetaBSON struct {
	req  RequestBSON
	res  ResponseBSON
	time time.Time
	body string
	from string
}

// RequestBSON represents a Request in Mongo (nested in Meta)
type RequestBSON struct{}

// ResponseBSON represents a Response in Mongo (nested in Meta)
type ResponseBSON struct{}

func (m *Meta) getType() Enum {
	if m.req != nil {
		return REQUEST
	} else if m.res != nil {
		return RESPONSE
	} else {
		return UNKNOWN
	}
}

// ToBSON transforms a Meta into a MetaBSON
func (m *Meta) ToBSON() (nr int64, err error) {
	if m.req != nil {
		buf, err2 := httputil.DumpRequest(m.req, false)
		if err2 != nil {
			return nr, err2
		}
		// write(&nr, &err, w, buf)
	} else if m.res != nil {
		buf, err2 := httputil.DumpResponse(m.resp, false)
		if err2 != nil {
			return nr, err2
		}
		// write(&nr, &err, w, buf)
	}

	return
}

// Save stores a Meta into the given Mongo collection
func (m *Meta) Save(coll *mgo.Collection) error {
	switch m.getType() {
	case REQUEST:
		coll.Insert(m.toBSON())
	case RESPONSE:
		coll.Update(m.id, m.toBSON())
	case REQUEST:
		coll.Insert(m.toBSON())
	}

}
