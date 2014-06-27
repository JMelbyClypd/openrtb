/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
 */
package apiserver

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"github.com/clyphub/tvapi/objects"
	"github.com/clyphub/tvapi/server"
	"github.com/clyphub/tvapi/store"
	"time"
)

const (
	OK int = 200
	CREATED int = 201
	ACCEPTED int = 202
	NOCONTENT int = 204
	BAD int = 400
	UNAUTH int = 401
	NOTFOUND int = 404
	ERROR int = 500
	PTV_VER string = `x-ptv-version`
	CT_JSON string = `application/json`
)


type ObjectProcessor interface {
	Unmarshal(buffer []byte) (obj objects.Storable, err error)
	ValidateRequest(msg objects.Storable) *server.CodedError
	ProcessRequest(req objects.Storable, responder *APIResponder) (resp objects.Storable, err *server.CodedError)
}


type APIResponder struct {
	server.BaseResponder
	method string
	path string
	store store.ObjectStore
	processor ObjectProcessor
}

func (r APIResponder) Init(srvr *server.Server, store store.ObjectStore) {
	r.store = store
	srvr.Register(r.method, r.path, r)
}

func (r APIResponder) Handle(req *http.Request) (status int, body []byte)  {

	log.Println("APIResponder.Handle entered")
	// Read the request body
	buffer, err := r.ReadBody(req)
	if (err != nil) {
		return BAD, nil
	}
	log.Println("APIResponder.Handle: body read")

	// Unmarshal the request
	obj, err := r.processor.Unmarshal(buffer)
	if (err != nil) {
		return ERROR, nil
	}
	log.Printf("Have object of type %T", obj)
	log.Printf("Received message with RequestId %s:\n%s\n", obj.GetKey(), buffer)

	// Validate it
	cerr := r.processor.ValidateRequest(obj)
	if(cerr != nil){
		log.Printf("Invalid request object: %s\n", cerr.Error())
		return cerr.Code(), nil
	}

	// Process it
	response, cerr := r.processor.ProcessRequest(obj, &r)
	if(cerr != nil){
		log.Printf("Error processing request object: %s\n", err.Error())
		return cerr.Code(), nil
	}

	// Marshal the results into the response body
	if(response != nil){
		buffer, err = objects.Marshal(response)
	}
	if (err != nil) {

		return ERROR, nil
	}
	log.Println("Request handled")
	return OK, buffer
}

type Callbacker struct {
	url string
	client *http.Client
}

func NewCallbacker(url string) *Callbacker{
	return &Callbacker{url, &http.Client{}}
}

func (cb Callbacker) Callback(b []byte) error {
	buffer := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", cb.url, buffer)
	if(err != nil){
		return err
	}
	req.Header.Set("Content-Type",CT_JSON)
	log.Printf("Sending message:\n%s", buffer)
	resp, err := cb.client.Do(req)
	if(err != nil){
		return err
	}
	if(resp.StatusCode != OK){
		return fmt.Errorf("Non-OK (%d) status received when sending callback to %s", resp.StatusCode, cb.url)
	}
	return nil
}

func (r APIResponder) waitAndSendResult(resp interface{}, url string, key string, delay int){
	log.Printf("waitAndSendResult called with url %s and key %s", url, key)

	// Let's go to sleep for a while
	time.Sleep(time.Duration(delay) * time.Second)

	buffer, err := objects.Marshal(resp)
	if(err != nil){
		log.Printf("Could not marshal callback object %s", key)
		return
	}
	cb := NewCallbacker(url)
	err = cb.Callback(buffer)
	if(err != nil){
		log.Printf("Could not send callback %s", key)
		return
	}
}


