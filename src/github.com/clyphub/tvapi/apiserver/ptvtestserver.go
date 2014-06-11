/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: main package for Programmatic TV API test service

 */
package main

import (
	"flag"
	"github.com/clyphub/tvapi/apiserver/apiserver"
	"github.com/clyphub/tvapi/server"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "addr", "127.0.0.1:12345", "bind host:port")
	flag.Parse()
}

func main() {
	s := server.NewServer()
	s.Init()
	s.AddResponder(apiserver.NewOrderAPIResponder())
	s.AddResponder(apiserver.NewInventoryAPIResponder())
	s.Open(address)
}
