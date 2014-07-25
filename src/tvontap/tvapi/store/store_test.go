package store

import (
	"testing"
	"tvontap/tvapi/objects"
)

var KEYS = []string{"orders", "orders", "buyer123", "reqabcd"}

func TestIsMatch(t *testing.T) {
	s := NewMapStore()
	check1 := []string{"orders", "orders", "buyer123"}
	if !s.isMatch(check1, KEYS) {
		t.Error("isMatch didn't return true")
		return
	}
	check1 = []string{"orders", "orders", "buyer124"}
	if s.isMatch(check1, KEYS) {
		t.Error("isMatch didn't return false")
		return
	}
}

func TestSaveGetErase(t *testing.T) {
	// Create the first test object
	obj := objects.AvailabilityRequestObject{}
	obj.RequestId = KEYS[3]
	obj.AdvertiserId = "Ducati"
	// Create the store
	ms := NewMapStore()
	// Save the first object
	e := ms.Set(KEYS, obj)
	if e != nil {
		t.Error(e)
	}
	// Retrieve and check the first object
	objg, e2 := ms.Get(KEYS)
	if e2 != nil {
		t.Error(e2)
	}
	at := objg.(objects.AvailabilityRequestObject).AdvertiserId
	if at != "Ducati" {
		t.Errorf("Wrong AdvertiserId value returned; got %s, expected %s", at, "Ducati")
	}
	// Create a second test object
	obj2 := objects.AvailabilityRequestObject{}
	obj2.RequestId = "abc999"
	obj2.AdvertiserId = "Aprilia"
	keys2 := KEYS
	keys2[3] = "abc999"
	// Save it
	e = ms.Set(keys2, obj2)
	if e != nil {
		t.Error(e)
	}
	// Retrieve both with a GetAll
	checkKeys := []string{"orders", "orders", "buyer123"}
	objArr, e := ms.GetAll(checkKeys)
	if len(objArr) != 2 {
		t.Errorf("Retrieved incorrect number of objects - expected %d got %d", 2, len(objArr))
	}
	// Erase the objects, then check to make sure at least one is gone
	e = ms.EraseAll(checkKeys)
	if e != nil {
		t.Error(e)
	}
	objg, e = ms.Get(KEYS)
	if objg != nil {
		t.Error("Erase didn't remove object")
	}
}
