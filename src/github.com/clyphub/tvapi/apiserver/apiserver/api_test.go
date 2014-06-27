/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Tests for apiserver package
*/
package apiserver


import (
	"log"
"testing"
	"time"
	"github.com/clyphub/tvapi/client"
	"github.com/clyphub/tvapi/objects"
"github.com/clyphub/tvapi/server"
	"github.com/clyphub/tvapi/store"
)

const (
	ADDR = "127.0.0.1:2345"
	CB_ADDR = "127.0.0.1:2345"
	CB_PATH = "/responses"
)

type CallbackReceiver struct {
	APIResponder
}

type CallbackProcessor struct {
	t *testing.T
	received objects.Storable
}

func NewCallbackReceiver(ts *testing.T) *CallbackReceiver {
	return &CallbackReceiver{APIResponder: APIResponder{path: CB_PATH, method: "POST", processor: &CallbackProcessor{t:ts}}}
}

func (r CallbackProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.AvailabilityResponseObject{}
	log.Printf("InventoryAPIProcessor unmarshalling with empty object of type %T", obj)
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r CallbackProcessor) ValidateRequest(rfp objects.Storable) *server.CodedError {
	return nil
}

func (r CallbackProcessor) ProcessRequest(rfp objects.Storable, responder *APIResponder) (resp objects.Storable, e *server.CodedError) {
	log.Println("Processing callback")
	r.received = rfp
	if(r.received == nil){
		log.Println("CallbackProcessor.received = nil")
		r.t.Fail()
	}
	return nil,nil
}

func (r CallbackReceiver) Init(srvr *server.Server, store store.ObjectStore) {
	srvr.Register("POST", CB_PATH, r)
}

func TestAvailabilityRequest(t *testing.T) {
	// Set up test server
	s := server.NewServer()
	s.Init()
	s.AddResponder(NewInventoryAPIResponder())
	cbr := NewCallbackReceiver(t)
	s.AddResponder(cbr)
	defer s.Close()
	go s.Open(ADDR)

	// Set up the client and test request
	req := objects.AvailabilityRequestObject{RequestId:"1234abc",BuyerId:"AcmeDSP123",AdvertiserId:"Ronco",ResponseUrl:"http://" + CB_ADDR + CB_PATH}
	cl, e := client.NewClient(ADDR)
	if(e != nil){
		t.Fatalf("Could not open client, error=%s",e.Error())
	}

	// Have the client do something useful
	e = cl.PostRequest(req, RFPPATH)
	if(e != nil){
		t.Logf("Transaction failed, error=%s", e.Error())
		t.Fail()
	}
	// This is an utter hack
	time.Sleep(time.Duration(5) * time.Second)

	t.Log("Done")

}
