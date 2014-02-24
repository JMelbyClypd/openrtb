package apiserver

import (
	"net/http"
	"testing"
)

func TestDispatcherGet(t *testing.T){
	var code int
	var buffer []byte
	obj := new(BaseResponder)
	req := new(http.Request)
	req.Method = "GET"
	code, buffer = obj.Handle(req)
	if(code != http.StatusOK){
		t.Errorf("Wrong code value returned; got %d, expected %d", code, http.StatusOK)
	}
	if( buffer != nil){
		t.Errorf("Wrong buffer value returned; got %s, expected nil", buffer)
	}
}

func TestMsgResponderGet(t *testing.T){
	var code int
	var buffer []byte
	obj := new(MsgResponder)
	obj.s = MSG
	req := new(http.Request)
	req.Method = "GET"
	code, buffer = obj.Handle(req)
	var msg = string(buffer)
	if(code != http.StatusOK){
		t.Errorf("Wrong code value returned; got %d, expected %d", code, http.StatusOK)
	}

	if( msg != MSG){
		t.Errorf("Wrong buffer value returned; got %s, expected %s", msg, MSG)
	}
}
