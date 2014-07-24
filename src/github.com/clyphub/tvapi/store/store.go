package store

import (
	"errors"
	"log"
	"github.com/clyphub/tvapi/objects"
	"github.com/clyphub/tvapi/util"
)

type ObjectStore interface {
	Set(keys []string, obj objects.Storable) error
	Get(keys []string) (objects.Storable, error)
	GetAll(keys []string) ([]objects.Storable, error)
	Erase(keys []string) error
	EraseAll(keys []string) error
}

type MapStore struct {
	store map[string]objects.Storable
}

func NewMapStore() MapStore {
	return MapStore{store: make(map[string]objects.Storable,10)}
}

func(s MapStore) isMatch(queryKeys []string, mapKeys []string) bool {
	lqk := len(queryKeys)
	if(lqk > len(mapKeys)) {
		return false
	}
	for i := 0;i<lqk;i++ {
		if(queryKeys[i] != mapKeys[i]){
			return false
		}
	}
	log.Printf("isMatch=true for %s and %s", util.Munge(queryKeys), util.Munge(mapKeys))
	return true
}

func(s MapStore) Set(keys []string, obj objects.Storable) error {
	if(obj == nil) {
		return errors.New("Object is nil")
	}
	if(len(keys) == 0) {
		return errors.New("Key is nil")
	}
	k := util.Munge(keys)
	s.store[k] = obj
	return nil
}

func(s MapStore) Get(keys []string) (objects.Storable, error){
	if(len(keys) == 0) {
		return nil, errors.New("No keys")
	}
	key := util.Munge(keys)
	return s.store[key], nil
}

func(s MapStore) GetAll(keys []string) ([]objects.Storable, error){
	if(len(keys) == 0) {
		return nil, errors.New("No keys")
	}
	ret := make([]objects.Storable,0,10)

	for key, value := range s.store {
		toks := util.Unmunge(key)
		if(s.isMatch(keys, toks)){
			ret = append(ret, value)
		}
	}
	return ret, nil
}

func(s MapStore) Erase(keys []string) error{
	if(len(keys)==0) {
		return errors.New("No keys")
	}
	key := util.Munge(keys)
	delete(s.store, key)
	return nil
}

func(s MapStore) EraseAll(keys []string) error{
	if(len(keys) == 0) {
		return errors.New("No keys")
	}
	for key, _ := range s.store {
		toks := util.Unmunge(key)
		if(s.isMatch(keys, toks)){
			s.Erase(toks)
		}
	}
	return nil
}

