/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
*/
package impl

import (
	"log"
	"net/http"
	"tvontap/tvapi/objects"
	"tvontap/tvapi/store"
	"tvontap/tvapi/util"
)

var RFPPATH = "/orders/availability/"

type InventoryRequestResponder struct {
	APIResponder
}

type InventoryRequestProcessor struct {
	StoreManager
}

func (r InventoryRequestProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.AvailabilityRequestObject{}
	log.Printf("InventoryAPIProcessor unmarshalling with empty object of type %T", obj)
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r InventoryRequestProcessor) ValidateRequest(pathTokens []string, queryTokens []string, rfp objects.Storable) *objects.CodedError {
	robj := rfp.(objects.AvailabilityRequestObject)
	if len(robj.RequestId) == 0 {
		return objects.NewError("No requestId", http.StatusBadRequest)
	}
	if len(robj.BuyerId) == 0 {
		return objects.NewError("No buyerId", http.StatusBadRequest)
	}
	if len(robj.AdvertiserId) == 0 {
		return objects.NewError("No advertiserId", http.StatusBadRequest)
	}
	if len(robj.ResponseUrl) == 0 {
		return objects.NewError("No responseUrl", http.StatusBadRequest)
	}
	log.Println("Request validated")
	return nil
}

func (r InventoryRequestProcessor) ProcessRequest(pathTokens []string, queryTokens []string, rfp objects.Storable, responder *APIResponder) ([]objects.Storable, *objects.CodedError) {
	log.Println("processRequest")
	// Save the request
	robj := rfp.(objects.AvailabilityRequestObject)
	err := r.SaveObject(robj.BuyerId, robj.RequestId, &robj)
	if err != nil {
		return nil, objects.NewError(err.Error(), http.StatusInternalServerError)
	}
	// Build the response
	respObj := &objects.AvailabilityResponseObject{}
	respObj.RequestId = robj.RequestId
	respObj.BuyerId = robj.BuyerId
	respObj.MinImpressions = robj.MinImpressions
	respObj.MaxImpressions = robj.MaxImpressions
	respObj.MinCPM = 15.30

	respObjs := make([]objects.Storable, 1)
	respObjs[0] = respObj

	// Send it later
	go responder.waitAndSendResult(respObjs, robj.ResponseUrl, robj.GetKey(), 1)

	return nil, nil
}

func NewInventoryRequestResponder(st store.ObjectStore) *InventoryRequestResponder {
	x := &InventoryRequestResponder{APIResponder{path: RFPPATH, processorMap: make(map[string]RequestProcessor, 4)}}
	x.AddProcessor("POST", &InventoryRequestProcessor{StoreManager{store: st, pathKeys: util.Unmunge(RFPPATH)}})
	x.AddProcessor("GET", &GenericGetProcessor{StoreManager{store: st, pathKeys: util.Unmunge(RFPPATH)}})
	x.AddProcessor("DELETE", &GenericDeleteProcessor{StoreManager{store: st, pathKeys: util.Unmunge(RFPPATH)}})
	return x
}
