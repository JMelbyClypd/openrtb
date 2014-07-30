#!/bin/sh

export GOPATH=`pwd`
go test -cover tvontap/tvapi/objects
go test -cover tvontap/tvapi/store
go test -cover tvontap/tvapi/util
go test -cover tvontap/tvapi/server
go test -cover tvontap/tvapi/apiserver/impl
