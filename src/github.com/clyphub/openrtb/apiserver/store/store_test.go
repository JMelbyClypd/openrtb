package store

import (
	"github.com/clyphub/openrtb/apiserver/objects"
	"testing"
)

const IDSTR = "abcd"

func TestGetField(t *testing.T){

	obj := new(objects.BidRequestObject)
	obj.Id = IDSTR
	obj.AuctionType = 8
	ms := NewMapStore()
	valstr,e := ms.getField("Id", obj)
	if(e != nil){
		t.Error(e)
	}
	if(valstr != IDSTR){
		t.Errorf("Wrong Id value returned; got %s, expected %s", valstr, IDSTR)
	}
	valstr,e = ms.getField("AuctionType", obj)
	if(e != nil){
		t.Error(e)
	}
	if(valstr != "8"){
		t.Errorf("Wrong AuctionType value returned; got %s, expected %s", valstr, "8")
	}
}

func TestSaveGetErase(t *testing.T){
	obj := new(objects.BidRequestObject)
	obj.Id = IDSTR
	obj.AuctionType = 8
	ms := NewMapStore()
	k,e := ms.Save(obj)
	if(e != nil){
		t.Error(e)
	}
	if(k != IDSTR){
		t.Errorf("Wrong Id value returned; got %s, expected %s", k, IDSTR)
	}
	obj2,e2 := ms.Get(IDSTR)
	if(e2 != nil){
		t.Error(e2)
	}
	at := obj2.(*objects.BidRequestObject).AuctionType
	if(at != 8){
		t.Errorf("Wrong AuctionType value returned; got %d, expected %d", at, 8)
	}
	_,e = ms.Erase(IDSTR)
	if(e != nil){
		t.Error(e)
	}
	obj2,e = ms.Get(IDSTR)
	if(obj2 != nil){
		t.Error("Erase didn't remove object")
	}
}
