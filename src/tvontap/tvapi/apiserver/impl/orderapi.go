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

var ORDERPATH = "/orders/order/"

type OrderAPIResponder struct {
	APIResponder
}

type OrderAPIProcessor struct {
	StoreManager
}

func (r OrderAPIProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.OrderObject{}
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r OrderAPIProcessor) ValidateRequest(pathTokens []string, queryTokens []string, order objects.Storable) *objects.CodedError {
	robj := order.(objects.OrderObject)
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

func (r OrderAPIProcessor) ProcessRequest(pathTokens []string, queryTokens []string, order objects.Storable, responder *APIResponder) ([]objects.Storable, *objects.CodedError) {
	log.Println("processRequest")
	robj := order.(objects.OrderObject)
	err := r.SaveObject(robj.BuyerId, robj.RequestId, &robj)
	if err != nil {
		return nil, objects.NewError(err.Error(), http.StatusInternalServerError)
	}
	// Build the response
	respObjs := make([]objects.Storable, 1)
	respObjs[0] = &objects.OrderAcceptanceObject{}

	// Send it later
	go responder.waitAndSendResult(respObjs, robj.ResponseUrl, robj.GetKey(), 1)
	return nil, nil
}

func NewOrderAPIResponder(st store.ObjectStore) *OrderAPIResponder {
	x := &OrderAPIResponder{APIResponder{path: ORDERPATH, processorMap: make(map[string]RequestProcessor, 4)}}
	x.AddProcessor("POST", &OrderAPIProcessor{StoreManager{store: st, pathKeys: util.Unmunge(ORDERPATH)}})
	x.AddProcessor("GET", &GenericGetProcessor{StoreManager{store: st, pathKeys: util.Unmunge(ORDERPATH)}})
	x.AddProcessor("DELETE", &GenericDeleteProcessor{StoreManager{store: st, pathKeys: util.Unmunge(ORDERPATH)}})
	return x
}
