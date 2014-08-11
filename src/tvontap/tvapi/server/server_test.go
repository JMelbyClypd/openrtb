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
	"strconv"
	"testing"
	"time"
	"tvontap/tvapi/store"
)

const (
	ADDR = "localhost"
	PORT = 2345
)

/*****************************************
Test the test dispatchers
*/

func _TestDispatcherGet(t *testing.T) {
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

func makeAddress(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

type Messager struct {
	BaseResponder
}

func (r *Messager) Init(srvr *Server) {
	srvr.Register("GET", "/msg/", r)
}

func (r *Messager) Handle(req *http.Request) (int, []byte) {
	log.Println("Messager.Handle")
	return http.StatusOK, []byte(MSG)
}

func TestServerGet(t *testing.T) {
	s := NewServer(store.NewMapStore())
	s.Init()
	s.AddResponder(&Messager{})
	defer s.Close()

	a := makeAddress(ADDR, PORT+4)
	go s.Open(a)

	time.Sleep(time.Duration(1) * time.Second)
	t.Log("Sending GET")
	DoGet(t, a, "/msg/", 200, MSG,
		func(r string, v string) bool {
			if r != v {
				t.Fail()
				t.Logf("Failed.  Expected %s got %s", r, v)
				return false
			}
			return true
		})

	t.Log("Done")
}

type Looper struct {
	BaseResponder
}

func (r *Looper) Init(srvr *Server) {
	srvr.Register("POST", "/loop", r)
	srvr.Register("PUT", "/loop", r)
	log.Println("Looper initialized")
}

func (r *Looper) Handle(req *http.Request) (int, []byte) {
	log.Println("Looper.Handle")
	buffer, e := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if e != nil {
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, buffer
}

func TestServerPost(t *testing.T) {
	s := NewServer(store.NewMapStore())
	s.Init()
	s.AddResponder(&Looper{})
	defer s.Close()

	a := makeAddress(ADDR, PORT+1)
	go s.Open(a)

	t.Log("Sending POST")
	DoPost(t, a, "/loop", "test", "application/text", 200, "test",
		func(r string, v string) bool {
			if r != v {
				t.Fail()
				t.Logf("Failed.  Expected %s got %s", r, v)
				return false
			}
			return true
		})

	t.Log("Done")
}

func TestServerPut(t *testing.T) {
	s := NewServer(store.NewMapStore())
	s.Init()
	s.AddResponder(&Looper{})
	defer s.Close()

	a := makeAddress(ADDR, PORT+2)
	go s.Open(a)

	DoPut(t, a, "/loop", "test", "application/text", 200, "test",
		func(r string, v string) bool {
			if r != v {
				t.Fail()
				t.Logf("Failed.  Expected %s got %s", r, v)
				return false
			}
			return true
		})

	t.Log("Done")

}

type Deleter struct {
	BaseResponder
}

func (r *Deleter) Init(srvr *Server) {
	srvr.Register("DELETE", "/obj", r)
	log.Println("Deleter initialized")
}

func (r *Deleter) Handle(req *http.Request) (status int, body []byte) {
	log.Println("Deleter.Handle")

	return http.StatusOK, nil
}

func TestServerDelete(t *testing.T) {
	s := NewServer(store.NewMapStore())
	s.Init()
	s.AddResponder(&Deleter{})
	defer s.Close()

	a := makeAddress(ADDR, PORT+3)

	go s.Open(a)

	DoDelete(t, a, "/obj", 200)

	t.Log("Done")

}
