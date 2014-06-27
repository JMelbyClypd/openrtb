/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Tests for objects package
*/

package objects

import (
	"testing"
)

func TestMarshalling(t *testing.T){
	obj := AvailabilityRequestObject{RequestId:"abc123"}
	buffer, e := Marshal(obj)
	if(e != nil){
		t.Logf("Marshal() failed with error %s", e.Error())
		t.Fail()
		return
	}
	obj2 := AvailabilityRequestObject{}
	e = Unmarshal(&obj2, buffer)
	if(e != nil){
		t.Logf("Unmarshal() failed with error %s", e.Error())
		t.Fail()
		return
	}
	if(obj2.RequestId != "abc123"){
		t.Logf("Unmarshal() failed returning bad object %s", obj2.RequestId)
		t.Fail()
		return
	}
}
