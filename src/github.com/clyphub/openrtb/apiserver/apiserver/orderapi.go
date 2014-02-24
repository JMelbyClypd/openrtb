/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
 */
package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/clyphub/openrtb/apiserver/objects"
	"github.com/clyphub/openrtb/apiserver/store"
)

var ORDERPATH = "/order/order"
var RFPPATH = "/order/rfp"

type APIResponder struct {
	BaseResponder
	method string
	path string
	store *store.Store
}

func (r *APIResponder) Init(srvr *Server, store *store.Store) {
	r.store = store
	srvr.register(r.method, r.path, r)
}

func (r *APIResponder) getRequestObject() interface{} {
	return nil
}

func (r *APIResponder) validateRequest(msg interface{}) error {
	return nil
}

func (r *APIResponder) processRequest(req interface{}) (resp interface{}, err error){
	return nil,nil
}

type RfpAPIResponder struct {
	APIResponder
}

func NewRfpAPIResponder() *RfpAPIResponder {
	var r RfpAPIResponder
	r.path = RFPPATH
	r.method = "POST"
	return &r
}

type OrderAPIResponder struct {
	APIResponder
}

func NewOrderAPIResponder() *OrderAPIResponder {
	var r OrderAPIResponder
	r.path = ORDERPATH
	r.method = "POST"
	return &r
}

func (r *RfpAPIResponder) getRequestObject() interface{} {
	var obj objects.RFPObject
	return obj
}

func (r *OrderAPIResponder) getRequestObject() interface{} {
	var obj objects.OrderObject
	return obj
}


func (r *APIResponder) Handle(req *http.Request) (status int, body []byte)  {

	// Read the request body
	var buffer []byte
	var code int
	var e error
	buffer, code = r.readBody(req)
	if (code != http.StatusOK) {
		return code, nil
	}

	// Now let's unmarshal the actual API request and get some work done
	//	First we need to guess the type of object
	var response interface{}
	var obj = r.getRequestObject()
	e = json.Unmarshal(buffer, &obj)
	if (e != nil) {
		log.Printf("Error unmarshalling request body (%s) error=%s",buffer,e.Error())
		return http.StatusInternalServerError, nil
	}
	e = r.validateRequest(obj)
	if(e != nil){
		log.Printf("Invalid request object: %s\n", e.Error())
		return 400, nil
	}
	response, e = r.processRequest(obj)

	// Marshal the results into the response body
	code,buffer = r.marshal(response)
	return code, buffer
}

func (r *RfpAPIResponder) validateRequest(rfp interface{}) error {
	robj := rfp.(objects.RFPObject)
	if(len(robj.BuyerId) == 0){

	}
	return nil
}

func (r *OrderAPIResponder) validateRequest(order interface{}) error {
	return nil
}

func (r *RfpAPIResponder) processRequest(rfp interface{}) (resp interface{}, err error){
	var obj objects.ProposalObject
	return obj,nil
}


func (r *OrderAPIResponder) processRequest(order interface{}) (resp interface{}, err error){
	var obj objects.AcceptanceObject
	return obj,nil
}



