package server

import (
	"net/http"
	"testing"
	"tvontap/tvapi/util"
)

type DummyHandler struct {
}

func (d DummyHandler) Handle(req *http.Request) (status int, body []byte) {
	return 200, nil
}

func TestSplit(t *testing.T) {
	val := "/Orders/Availability/"

	arr := util.Unmunge(val)
	t.Logf("Split %s into: ", val)
	for _, tok := range arr {
		if len(tok) > 0 {
			t.Log(tok)
		}
	}
	if len(arr) != 2 {
		t.Fatalf("Wrong number, got %d", len(arr))
	}
}

func TestRouter(t *testing.T) {
	bmh := BadMethodHandler{}
	dh := &DummyHandler{}
	dh2 := &DummyHandler{}
	r := &Router{root: newNode("/"), badMethodHandler: bmh}
	r.Register("POST", "/orders/availability/", dh)
	r.Register("GET", "/orders/availability/", dh)
	r.Register("GET", "/orders/", dh2)
	r.dumpNodes()
	mh := r.resolveHandler("DELETE", "/orders/availability")
	if mh != bmh {
		t.Fatal("Got non-nil result when expecting nil (bogus method)")
	}
	mh = r.resolveHandler("POST", "/orders/bogus/")
	if mh != bmh {
		t.Fatal("Got non-nil result when expecting nil (bogus path")
	}
	mh = r.resolveHandler("POST", "/orders/availability/")
	if mh != dh {
		t.Fatal("Got nil result when expecting non-nil (on POST)")
	}
	mh = r.resolveHandler("GET", "/orders/availability/")
	if mh != dh {
		t.Fatal("Got nil result when expecting non-nil (GET availability)")
	}
	mh = r.resolveHandler("GET", "/orders/")
	if mh != dh2 {
		t.Fatal("Got nil result when expecting non-nil (GET orders)")
	}
	mh = r.resolveHandler("GET", "/orders/availability/buyer1234/order123/")
	if mh != dh {
		t.Fatal("Got nil result when expecting non-nil (GET full path)")
	}
}
