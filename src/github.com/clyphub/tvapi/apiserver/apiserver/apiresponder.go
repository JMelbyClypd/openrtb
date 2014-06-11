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

type APIResponder struct {
	server.BaseResponder
	method string
	path string
	store store.ObjectStore
}

func (r *APIResponder) Init(srvr *server.Server, store store.ObjectStore) {
	r.store = store
	srvr.Register(r.method, r.path, r)
}

func (r *APIResponder) ValidateRequest(msg interface{}) error {
	return nil
}

func (r *APIResponder) ProcessRequest(req interface{}) (resp interface{}, err error){
	return nil,nil
}

func (r *APIResponder) Handle(req *http.Request) (status int, body []byte)  {

	// Read the request body
	buffer, err := r.ReadBody(req)
	if (err != nil) {
		return BAD, nil
	}

	// Unmarshal the request
	obj := r.GetObject()
	err = objects.Unmarshal(obj,buffer)
	if (err != nil) {
		return ERROR, nil
	}

	// Validate it
	err = r.ValidateRequest(obj)
	if(err != nil){
		log.Printf("Invalid request object: %s\n", err.Error())
		return BAD, nil
	}

	// Process it
	response, err := r.ProcessRequest(obj)

	// Marshal the results into the response body
	if(response != nil){
		buffer, err = objects.Marshal(response)
	}
	if (err != nil) {
		return ERROR, nil
	}
	return OK, buffer
}

type Callbacker struct {
	url string
	client *http.Client
}

func NewCallbacker(url string) *Callbacker{
	return &Callbacker{url, &http.Client{}}
}

func (cb *Callbacker) Callback(b []byte) error {
	buffer := bytes.NewBuffer(b)
	req, err := http.NewRequest("PUT", cb.url, buffer)
	if(err != nil){
		return err
	}
	req.Header.Set("Content-Type",CT_JSON)
	resp, err := cb.client.Do(req)
	if(err != nil){
		return err
	}
	if(resp.StatusCode != OK){
		return fmt.Errorf("Non-OK (%d) status received when sending callback to %s", resp.StatusCode, cb.url)
	}
	return nil
}




