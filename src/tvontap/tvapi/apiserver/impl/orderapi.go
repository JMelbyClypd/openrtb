/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
*/
package impl

import (
	"log"
	"tvontap/tvapi/objects"
	"tvontap/tvapi/store"
	"tvontap/tvapi/util"
)

var ORDERPATH = "/orders/order"

type OrderAPIResponder struct {
	APIResponder
}

type OrderAPIProcessor struct {
	StoreManager
}

func NewOrderAPIResponder(st store.ObjectStore) *OrderAPIResponder {
	x := &OrderAPIResponder{APIResponder{path: ORDERPATH, processorMap: make(map[string]RequestProcessor, 4)}}
	x.AddProcessor("POST", &OrderAPIProcessor{StoreManager{store: st, pathKeys: util.Unmunge(ORDERPATH)}})
	return x
}

func (r OrderAPIProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.OrderObject{}
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r OrderAPIProcessor) ValidateRequest(pathTokens []string, queryTokens []string, order objects.Storable) *objects.CodedError {
	return nil
}

func (r OrderAPIProcessor) ProcessRequest(pathTokens []string, queryTokens []string, order objects.Storable, responder *APIResponder) ([]objects.Storable, *objects.CodedError) {
	log.Println("processRequest")
	respObjs := make([]objects.Storable, 1)
	respObjs[0] = &objects.OrderAcceptanceObject{}
	return respObjs, nil
}
