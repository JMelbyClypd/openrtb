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
)

var ORDERPATH = "/order/order"

type OrderAPIResponder struct {
	APIResponder
}

type OrderAPIProcessor struct {

}

func NewOrderAPIResponder() *OrderAPIResponder {
	return &OrderAPIResponder{APIResponder{path: ORDERPATH, method: "POST", processor:&OrderAPIProcessor{}}}
}

func (r OrderAPIProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.OrderObject{}
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r OrderAPIProcessor) ValidateRequest(order objects.Storable) *server.CodedError {
	return nil
}

func (r OrderAPIProcessor) ProcessRequest(order objects.Storable, responder *APIResponder) (objects.Storable, *server.CodedError){
	log.Println("processRequest")
	var obj objects.OrderAcceptanceObject
	return obj,nil
}
