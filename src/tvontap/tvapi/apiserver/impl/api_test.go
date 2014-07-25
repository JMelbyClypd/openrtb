/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Tests for apiserver package
*/
package impl

import (
	"log"
	"testing"
	"time"
	"tvontap/tvapi/client"
	"tvontap/tvapi/objects"
	"tvontap/tvapi/server"
	"tvontap/tvapi/store"
	"tvontap/tvapi/util"
)

const (
	ADDR    = "127.0.0.1:2345"
	CB_ADDR = "127.0.0.1:2345"
	CB_PATH = "/responses/"
)

type CallbackReceiver struct {
	APIResponder
}

type CallbackProcessor struct {
	t        *testing.T
	received objects.Storable
	StoreManager
}

func NewCallbackReceiver(ts *testing.T) *CallbackReceiver {
	x := &CallbackReceiver{APIResponder: APIResponder{path: CB_PATH, processorMap: make(map[string]RequestProcessor, 4)}}
	cp := CallbackProcessor{t: ts}
	cp.pathKeys = util.Unmunge(CB_PATH)
	x.AddProcessor("POST", &cp)
	return x
}

func (r CallbackProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	obj := objects.AvailabilityResponseObject{}
	log.Printf("InventoryAPIProcessor unmarshalling with empty object of type %T", obj)
	err := objects.Unmarshal(&obj, buffer)
	return obj, err
}

func (r CallbackProcessor) ValidateRequest(pathTokens []string, queryTokens []string, rfp objects.Storable) *objects.CodedError {
	if &rfp == nil {
		log.Println("CallbackProcessor.received = nil")
		r.t.Fail()
		return objects.NewError("No response received", 400)
	}
	return nil
}

func (r CallbackProcessor) ProcessRequest(pathTokens []string, queryTokens []string, rfp objects.Storable, responder *APIResponder) (resp []objects.Storable, e *objects.CodedError) {
	log.Println("Processing callback")
	r.received = rfp
	return nil, nil
}

func (r CallbackReceiver) Init(srvr *server.Server) {
	srvr.Register("POST", CB_PATH, r)
}

func TestAvailabilityRequest(t *testing.T) {
	// Set up test server
	log.Println("Setting up test store and test server")
	teststore := store.NewMapStore()
	testserver := server.NewServer(teststore)
	testserver.Init()
	testserver.AddResponder(NewInventoryRequestResponder(teststore))
	cbr := NewCallbackReceiver(t)
	testserver.AddResponder(cbr)
	defer testserver.Close()
	log.Println("Setup complete, opening server")
	go testserver.Open(ADDR)

	// Set up the client and test request
	req := objects.AvailabilityRequestObject{RequestId: "1234abc", BuyerId: "AcmeDSP123", AdvertiserId: "Ronco", ResponseUrl: "http://" + CB_ADDR + CB_PATH}
	cl, ex := client.NewClient(ADDR)
	if ex != nil {
		t.Fatalf("Could not open client, error=%s", ex.Error())
	}

	log.Println("Client set up")
	// Have the client do something useful
	err := cl.PostRequest(req, RFPPATH)
	if err != nil {
		t.Logf("Transaction failed, error=%s", err.Error())
		t.Fail()
	}
	log.Println("Sent POST request")
	// This is an utter hack
	time.Sleep(time.Duration(1) * time.Second)

	// Test GET
	igot := make([]objects.AvailabilityRequestObject, 1)
	err = cl.GetRequest("http://"+ADDR+RFPPATH+"AcmeDSP123/1234abc/", &igot)
	if err != nil {
		t.Logf("Get failed, error=%s", err.Error())
		t.Fail()
		return
	}
	log.Println("Requested GET")
	if len(igot) == 0 {
		t.Log("Get failed, got 0 responses")
		t.Fail()
		return
	}
	if igot[0].RequestId != "1234abc" {
		t.Logf("Get failed, got %s", igot[0].RequestId)
		t.Fail()
		return
	}

	// Test DELETE
	err = cl.DelRequest("http://"+ADDR+RFPPATH+"AcmeDSP123/1234abc/")
	if err != nil {
		t.Logf("Delete failed, error=%s", err.Error())
		t.Fail()
		return
	}
	log.Println("Requested DELETE")
	igot = make([]objects.AvailabilityRequestObject, 1)
	err = cl.GetRequest("http://"+ADDR+RFPPATH+"AcmeDSP123/1234abc/", &igot)
	if(err == nil || err.Code() != 404){
		t.Log("Get following Delete failed, should have returned a 404")
		t.Fail()
		return
	}
	t.Log("Done")

}
