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
	"strings"
	"testing"
	"github.com/clyphub/tvapi/store"
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

/*****************************************
	Helper functions
 */

func do(t *testing.T, mthd string, addr string, path string, body string, bodyType string, expected string, checker func(ref string, val string) bool, stopper func()) {
	url := "http://" + addr + path


	req, err := http.NewRequest(mthd, url, strings.NewReader(body))


	if(len(bodyType) != 0) {
		req.Header.Add("Content-Type", bodyType)
	}

	if err != nil {
		t.Log("Could not create request")
		t.Fail()
		stopper()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Log("Could not send request")
		t.Fail()
		stopper()
		return
	}

	status := resp.StatusCode
	if status != 200 {
		t.Logf("Bad response: %d", status)
		t.Fail()
		stopper()
		return
	}

	var buffer []byte
	buffer, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("Couldn't read body")
		t.Fail()
		stopper()
		return
	}
	if(len(buffer)>0){
		msg := string(buffer)
		var passes = true
		if(checker != nil){
			passes = checker(expected, msg)
		}

		if passes != true {
			t.Fail()
		} else {

		}
	} else {
		t.Log("OK")
	}



	stopper()
}

func DoGet(t *testing.T, addr string, path string, expected string, checker func(ref string, val string) bool, stopper func()) {
	do(t, "GET", addr, path, "", "", expected, checker, stopper)
}

func DoPost(t *testing.T, addr string, path string, body string, bodyType string, expected string, checker func(ref string, val string) bool, stopper func()) {
	do(t, "POST", addr, path, body, bodyType, expected, checker, stopper)
}

func DoPut(t *testing.T, addr string, path string, body string, bodyType string, expected string, checker func(ref string, val string) bool, stopper func()) {
	do(t, "PUT", addr, path, body, bodyType, expected, checker, stopper)
}

func DoDelete(t *testing.T, addr string, path string, stopper func()) {
	do(t, "DELETE", addr, path, "", "", "", nil, stopper)
}

/*****************************************
	Test the server
 */

func TestServerGet(t *testing.T) {
	s := NewServer()
	s.Init()
	defer s.Close()

	go DoGet(t, ADDR, "/",MSG,
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

func (r *Looper) Handle(req *http.Request) (status int, body []byte) {
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

	go DoPost(t, ADDR, "/loop","test","application/text", "test",
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

	go DoPut(t, ADDR, "/loop","test","application/text", "test",
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

	go DoDelete(t, ADDR, "/obj",
		func() {
			s.Close()
		})

	s.Open(ADDR)
	t.Log("Done")

}
