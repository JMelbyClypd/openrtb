/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
 */
package apiserver

import (
	"log"
	"github.com/clyphub/tvapi/objects"
	"github.com/clyphub/tvapi/server"
	"net/http"
)

var RFPPATH = "/order/availability"

type InventoryAPIResponder struct {
	APIResponder
}

type InventoryAPIProcessor struct {
}

func (r InventoryAPIProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.AvailabilityRequestObject{}
	log.Printf("InventoryAPIProcessor unmarshalling with empty object of type %T", obj)
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r InventoryAPIProcessor) ValidateRequest(rfp objects.Storable) *server.CodedError {
	robj := rfp.(objects.AvailabilityRequestObject)
	if(len(robj.RequestId) == 0){
		return server.NewError("No requestId", http.StatusBadRequest)
	}
	if(len(robj.BuyerId) == 0){
		return server.NewError("No buyerId", http.StatusBadRequest)
	}
	if(len(robj.AdvertiserId) == 0){
		return server.NewError("No advertiserId", http.StatusBadRequest)
	}
	if(len(robj.ResponseUrl) == 0){
		return server.NewError("No responseUrl", http.StatusBadRequest)
	}
	log.Println("Request validated")
	return nil
}

func (r InventoryAPIProcessor) ProcessRequest(rfp objects.Storable, responder *APIResponder) (objects.Storable, *server.CodedError){
	log.Println("processRequest")
	// Save the request
	robj := rfp.(objects.AvailabilityRequestObject)
	_, err := responder.store.Save(&robj)
	if(err != nil){
		return nil, server.NewError(err.Error(), http.StatusInternalServerError)
	}
	// Build the response
	respObj := &objects.AvailabilityResponseObject{}
	respObj.RequestId = robj.RequestId
	respObj.BuyerId = robj.BuyerId
	respObj.MinImpressions = robj.MinImpressions
	respObj.MaxImpressions = robj.MaxImpressions
	respObj.MinCPM = 15.30

	// Send it later
	go responder.waitAndSendResult(respObj, robj.ResponseUrl, robj.GetKey(), 1)

	return nil,nil
}

func NewInventoryAPIResponder() *InventoryAPIResponder {
	return &InventoryAPIResponder{APIResponder{path: RFPPATH, method: "POST", processor:&InventoryAPIProcessor{}}}
}
