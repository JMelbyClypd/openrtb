/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
 */
package apiserver

import (
	"log"
	"github.com/clyphub/tvapi/objects"
)

var ORDERPATH = "/order/order"

type OrderAPIResponder struct {
	APIResponder
}

func NewOrderAPIResponder() *OrderAPIResponder {
	var r OrderAPIResponder
	r.path = ORDERPATH
	r.method = "POST"
	return &r
}

func (r *OrderAPIResponder) getObject() interface{} {
	var obj objects.OrderObject
	return obj
}

func (r *OrderAPIResponder) validateRequest(order interface{}) error {
	return nil
}

func (r *OrderAPIResponder) processRequest(order interface{}) (resp interface{}, err error){
	log.Println("processRequest")
	var obj objects.OrderAcceptanceObject
	return obj,nil
}
