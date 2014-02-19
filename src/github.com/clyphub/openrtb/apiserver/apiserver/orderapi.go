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

type RfpAPIResponder struct {
	MethodDispatcher
}

func (r *RfpAPIResponder) Init(srvr *Server, resp *Router, store *store.Store) {
	resp.register(RFPPATH, r)
	r.sink = r
}

type OrderAPIResponder struct {
	MethodDispatcher
}

func (r *OrderAPIResponder) Init(srvr *Server, resp *Router, store *store.Store) {
	resp.register(ORDERPATH, r)
	r.sink = r
}

func (r *RfpAPIResponder) HandlePost(req *http.Request) (status int, body []byte)  {

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
	var obj objects.RFPObject
	e = json.Unmarshal(buffer, &obj)
	if (e != nil) {
		log.Printf("Error unmarshalling request body (%s) error=%s",buffer,e.Error())
		return http.StatusInternalServerError, nil
	}
	e = r.validateRfpRequest(obj)
	if(e != nil){
		log.Printf("Invalid RFP: %s\n", e.Error())
		return 400, nil
	}
	response, e = r.processRfpRequest(obj)

	// Marshal the results into the response body
	code,buffer = r.marshal(response)
	return code, buffer
}

func (r *OrderAPIResponder) HandlePost(req *http.Request) (status int, body []byte)  {

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
	var obj objects.OrderObject
	e = json.Unmarshal(buffer, &obj)
	if (e != nil) {
		log.Printf("Error unmarshalling request body (%s) error=%s",buffer,e.Error())
		return http.StatusInternalServerError, nil
	}
	e = r.validateOrderRequest(obj)
	if(e != nil){
		log.Printf("Invalid Order: %s\n", e.Error())
		return 400, nil
	}
	response, e = r.processOrderRequest(obj)

	// Marshal the results into the response body
	code,buffer = r.marshal(response)
	return code, buffer
}

func (r *RfpAPIResponder) validateRfpRequest(rfp objects.RFPObject) error {
	return nil
}

func (r *OrderAPIResponder) validateOrderRequest(order objects.OrderObject) error {
	return nil
}

func (r *RfpAPIResponder) processRfpRequest(rfp objects.RFPObject) (resp objects.ProposalObject, err error){
	var obj objects.ProposalObject
	return obj,nil
}


func (r *OrderAPIResponder) processOrderRequest(order objects.OrderObject) (resp objects.AcceptanceObject, err error){
	var obj objects.AcceptanceObject
	return obj,nil
}

func (r *OrderAPIResponder) HandleGet(req *http.Request) (status int, body []byte) {

	return http.StatusOK, nil
}

func (r *RfpAPIResponder) HandleGet(req *http.Request) (status int, body []byte) {
	return http.StatusOK, nil
}

func (r *RfpAPIResponder) HandlePut(req *http.Request) (status int, body []byte) {
	return http.StatusOK, nil
}

func (r *OrderAPIResponder) HandlePut(req *http.Request) (status int, body []byte) {
	return http.StatusOK, nil
}

