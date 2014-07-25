/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: main package for Programmatic TV API test service

*/
package main

import (
	"flag"
	"tvontap/tvapi/apiserver/impl"
	"tvontap/tvapi/server"
	"tvontap/tvapi/store"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "addr", "127.0.0.1:12345", "bind host:port")
	flag.Parse()
}

func main() {
	st := store.NewMapStore()
	s := server.NewServer(st)
	s.Init()
	s.AddResponder(impl.NewOrderAPIResponder(st))
	s.AddResponder(impl.NewInventoryRequestResponder(st))
	s.Open(address)
}
