package main

import (
	"flag"
	"github.com/clyphub/openrtb/apiserver/apiserver"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "addr", "127.0.0.1:12345", "bind host:port")
	flag.Parse()
}

func main() {
	server := new(apiserver.Server)
	server.Init()
	server.AddResponder(new(apiserver.OrderAPIResponder))
	server.AddResponder(new(apiserver.RfpAPIResponder))
	server.Open(address)
}
