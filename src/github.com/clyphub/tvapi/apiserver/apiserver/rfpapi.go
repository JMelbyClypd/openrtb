/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
 */
package apiserver

import (
	"log"
	"github.com/clyphub/tvapi/objects"
	"time"
)

var RFPPATH = "/order/availability"

type InventoryAPIResponder struct {
	APIResponder
}

func NewInventoryAPIResponder() *InventoryAPIResponder {
	return &InventoryAPIResponder{APIResponder{path: RFPPATH, method: "POST"}}
}

func (r *InventoryAPIResponder) GetObject() interface{} {
	return objects.AvailabilityRequestObject{}
}

func (r *InventoryAPIResponder) validateRequest(rfp interface{}) error {
	robj := rfp.(objects.AvailabilityRequestObject)
	if(len(robj.BuyerId) == 0){

	}
	return nil
}

func (r *InventoryAPIResponder) processRequest(rfp interface{}) (resp interface{}, e error){
	log.Println("processRequest")
	robj := rfp.(objects.AvailabilityRequestObject)
	_, e = r.store.Save(&robj)
	if(e != nil){
		return nil, e
	}
	go r.waitAndSendResult(robj)

	return nil,nil
}

func (r *InventoryAPIResponder) waitAndSendResult(obj objects.AvailabilityRequestObject){
	// Let's go to sleep for a while
	time.Sleep(100 * time.Second)

	// Extract the response URL
	url := obj.ResponseUrl
	key := obj.GetKey()

	resp := &objects.AvailabilityResponseObject{}
	buffer, err := objects.Marshal(resp)
	if(err != nil){
		log.Printf("Could not marshal availability response %s", key)
		return
	}
	cb := NewCallbacker(url)
	err = cb.Callback(buffer)
	if(err != nil){
		log.Printf("Could not send availability response %s", key)
		return
	}
}
