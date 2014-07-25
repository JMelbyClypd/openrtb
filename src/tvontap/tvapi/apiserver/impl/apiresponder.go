/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
*/
package impl

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"tvontap/tvapi/objects"
	"tvontap/tvapi/server"
	"tvontap/tvapi/store"
	"tvontap/tvapi/util"
)

const (
	OK        int    = 200
	CREATED   int    = 201
	ACCEPTED  int    = 202
	NOCONTENT int    = 204
	BAD       int    = 400
	UNAUTH    int    = 401
	NOTFOUND  int    = 404
	ERROR     int    = 500
	PTV_VER   string = `x-ptv-version`
	CT_JSON   string = `application/json`
)

type RequestProcessor interface {
	Unmarshal(buffer []byte) (objects.Storable, error)
	ValidateRequest(pathTokens []string, queryTokens []string, msg objects.Storable) *objects.CodedError
	ProcessRequest(pathTokens []string, queryTokens []string, req objects.Storable, responder *APIResponder) ([]objects.Storable, *objects.CodedError)
}

type StoreManager struct {
	store    store.ObjectStore
	pathKeys []string
}

func (sm StoreManager) SaveObject(buyerId string, objectKey string, obj objects.Storable) error {
	log.Printf("SaveObject(), buyerId=%s, objectKey=%s", buyerId, objectKey)
	keys := append(sm.pathKeys, buyerId, objectKey)
	if sm.store == nil {
		log.Fatalln("Store is nil")
	}
	return sm.store.Set(keys, obj)
}

func (sm StoreManager) DeleteObject(buyerId string, objectKey string) error {
	keys := append(sm.pathKeys, buyerId, objectKey)
	return sm.store.Erase(keys)
}

func (sm StoreManager) DeleteAllByBuyer(buyerId string) error {
	keys := append(sm.pathKeys, buyerId)
	return sm.store.EraseAll(keys)
}

func (sm StoreManager) GetObject(buyerId string, objectKey string) (objects.Storable, error) {
	keys := append(sm.pathKeys, buyerId, objectKey)
	return sm.store.Get(keys)
}

func (sm StoreManager) GetAllByBuyer(buyerId string) ([]objects.Storable, error) {
	keys := append(sm.pathKeys, buyerId)
	return sm.store.GetAll(keys)
}

type APIResponder struct {
	server.BaseResponder
	path         string
	processorMap map[string]RequestProcessor
}

func (r APIResponder) Init(srvr *server.Server) {
	log.Printf("Init() called for responder %s", r.path)
	for key, _ := range r.processorMap {
		srvr.Register(key, r.path, r)
	}
}

func (r APIResponder) AddProcessor(meth string, rp RequestProcessor) {
	r.processorMap[meth] = rp
}

func (r APIResponder) getProcessor(method string) RequestProcessor {
	if &method == nil {
		log.Println("no method value provided to getProcessor()")
		return nil
	}
	return r.processorMap[method]
}

func (r APIResponder) Handle(req *http.Request) (int, []byte) {

	log.Println("APIResponder.Handle entered")
	method := req.Method
	u := req.URL
	pathTokens := util.Unmunge(u.Path)
	queryTokens := strings.Split(u.RawQuery, "&")

	// Read the request body
	buffer, err := r.ReadBody(req)
	if err != nil {
		return BAD, nil
	}
	log.Println("APIResponder.Handle: body read")

	pr := r.getProcessor(method)
	if &pr == nil {
		return BAD, nil
	}

	obj, err := pr.Unmarshal(buffer)
	if err != nil {
		return ERROR, nil
	}
	if obj != nil {
		log.Printf("Have object of type %T", obj)
		log.Printf("Received message with RequestId %s:\n%s\n", obj.GetKey(), buffer)
	} else {
		log.Println("Request did not include a body")
	}
	cerr := pr.ValidateRequest(pathTokens, queryTokens, obj)
	if cerr != nil {
		log.Printf("Invalid %s request: %s\n", method, cerr.Error())
		return cerr.Code(), nil
	}

	// Process it
	responses, cerr := pr.ProcessRequest(pathTokens, queryTokens, obj, &r)
	if cerr != nil {
		log.Printf("Error processing %s request: %s\n", method, cerr.Error())
		return cerr.Code(), nil
	}

	// Marshal the results into the response body
	if len(responses) > 0 {
		buffer, err = objects.Marshal(responses)
		log.Printf("Marshalling %d objects into response as \n%s", len(responses), buffer)
	}
	if err != nil {
		log.Printf("Error marshalling response: %s\n", err.Error())
		return ERROR, nil
	}
	log.Printf("%s Request handled", method)
	return OK, buffer

}

type Callbacker struct {
	url    string
	client *http.Client
}

func NewCallbacker(url string) *Callbacker {
	return &Callbacker{url, &http.Client{}}
}

func (cb Callbacker) Callback(b []byte) error {
	buffer := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", cb.url, buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", CT_JSON)
	log.Printf("Sending message:\n%s", buffer)
	resp, err := cb.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != OK {
		return fmt.Errorf("Non-OK (%d) status received when sending callback to %s", resp.StatusCode, cb.url)
	}
	return nil
}

func (r APIResponder) waitAndSendResult(resp []objects.Storable, url string, key string, delay int) {
	log.Printf("waitAndSendResult called with url %s and key %s", url, key)

	// Let's go to sleep for a while
	time.Sleep(time.Duration(delay) * time.Second)
	var buffer []byte
	var err error

	if len(resp) == 1 {
		buffer, err = objects.Marshal(resp[0])
	} else {
		buffer, err = objects.Marshal(resp)
	}

	if err != nil {
		log.Printf("Could not marshal callback object %s", key)
		return
	}
	cb := NewCallbacker(url)
	err = cb.Callback(buffer)
	if err != nil {
		log.Printf("Could not send callback %s", key)
		return
	}
}
