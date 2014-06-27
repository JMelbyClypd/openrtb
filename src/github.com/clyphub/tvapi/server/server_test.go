/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Tests for server package
*/
package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"github.com/clyphub/tvapi/store"
	"github.com/clyphub/tvapi/testutil"
)

const (
	ADDR = "127.0.0.1:2345"
)

/*****************************************
	Test the test dispatchers
 */

func TestDispatcherGet(t *testing.T) {
	var code int
	var buffer []byte
	obj := new(BaseResponder)
	req := new(http.Request)
	req.Method = "GET"
	code, buffer = obj.Handle(req)
	if code != http.StatusOK {
		t.Errorf("Wrong code value returned; got %d, expected %d", code, http.StatusOK)
	}
	if buffer != nil {
		t.Errorf("Wrong buffer value returned; got %s, expected nil", buffer)
	}
}

func TestMsgResponderGet(t *testing.T) {
	var code int
	var buffer []byte
	obj := new(MsgResponder)
	obj.s = MSG
	req := new(http.Request)
	req.Method = "GET"
	code, buffer = obj.Handle(req)
	var msg = string(buffer)
	if code != http.StatusOK {
		t.Errorf("Wrong code value returned; got %d, expected %d", code, http.StatusOK)
	}

	if msg != MSG {
		t.Errorf("Wrong buffer value returned; got %s, expected %s", msg, MSG)
	}
}

func TestServerGet(t *testing.T) {
	s := NewServer()
	s.Init()
	defer s.Close()

	go testutil.DoGet(t, ADDR, "/",200, MSG,
		func(r string, v string) bool {
			if r != v {
				t.Fail()
				t.Logf("Failed.  Expected %s got %s", r, v)
				return false
			}
			return true
		},
		func() {
			s.Close()
		})

	s.Open(ADDR)
	t.Log("Done")

}


type Looper struct {
	BaseResponder
}

func (r *Looper) Init(srvr *Server, store store.ObjectStore) {
	srvr.Register("POST", "/loop", r)
	srvr.Register("PUT", "/loop", r)
	log.Println("Looper initialized")
}

func (r *Looper) Handle(req *http.Request) (int, []byte) {
	log.Println("Looper.Handle")
	buffer, e := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if(e != nil){
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, buffer
}

func TestServerPost(t *testing.T) {
	s := NewServer()
	s.Init()
	s.AddResponder(&Looper{})
	defer s.Close()

	go testutil.DoPost(t, ADDR, "/loop","test","application/text", 200, "test",
		func(r string, v string) bool {
			if r != v {
				t.Fail()
				t.Logf("Failed.  Expected %s got %s", r, v)
				return false
			}
			return true
		},
		func() {
			s.Close()
		})

	s.Open(ADDR)
	t.Log("Done")

}

func TestServerPut(t *testing.T) {
	s := NewServer()
	s.Init()
	s.AddResponder(&Looper{})
	defer s.Close()

	go testutil.DoPut(t, ADDR, "/loop","test","application/text", 200, "test",
		func(r string, v string) bool {
			if r != v {
				t.Fail()
				t.Logf("Failed.  Expected %s got %s", r, v)
				return false
			}
			return true
		},
		func() {
			s.Close()
		})

	s.Open(ADDR)
	t.Log("Done")

}

type Deleter struct {
	BaseResponder
}

func (r *Deleter) Init(srvr *Server, store store.ObjectStore) {
	srvr.Register("DELETE", "/obj", r)
	log.Println("Deleter initialized")
}

func (r *Deleter) Handle(req *http.Request) (status int, body []byte) {
	log.Println("Deleter.Handle")

	return http.StatusOK, nil
}

func TestServerDelete(t *testing.T) {
	s := NewServer()
	s.Init()
	s.AddResponder(&Deleter{})
	defer s.Close()

	go testutil.DoDelete(t, ADDR, "/obj", 200,
		func() {
			s.Close()
		})

	s.Open(ADDR)
	t.Log("Done")

}
