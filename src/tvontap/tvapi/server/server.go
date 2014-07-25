/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Lightweight container for Programmatic TV API service
*/
package server

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"tvontap/tvapi/store"
)

const (
	MSG      = "Hey dere, dis is Hibbing callin\n"
	BYE_PATH = "/bye/"
	BYE_MSG  = "Buh-bye"
)

type Server struct {
	Router
	listener net.Listener
	store    store.ObjectStore
}

func NewServer(st store.ObjectStore) *Server {
	s := Server{}
	s.root = newNode("/")
	s.badMethodHandler = BadMethodHandler{}
	s.store = st
	//return &Server{Router{root:newNode("/"), badMethodHandler:BadMethodHandler{}}}
	return &s
}

type MethodHandler interface {
	Handle(req *http.Request) (status int, body []byte)
}

type Responder interface {
	MethodHandler
	Init(srvr *Server)
}

func (s Server) Init() {
	http.DefaultServeMux = http.NewServeMux()
	http.DefaultServeMux.Handle("/", s)

	mr := &MsgResponder{}
	mr.Init(&s)

	sr := &Downer{}
	sr.Init(&s)
}

func (s Server) Open(laddr string) {
	log.Println("Server opening")
	var err error
	s.listener, err = net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}
	http.Serve(s.listener, nil)

}

func (s Server) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
}

func (srvr Server) AddResponder(responder Responder) {
	responder.Init(&srvr)
}

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	log.Printf("Router.ServerHTTP entered with method %s and path %s", method, path)

	status := http.StatusBadRequest
	var body []byte

	sink := s.resolveHandler(method, path)
	status, body = sink.Handle(req)

	log.Printf("ServeHTTP: Writing response with status %d\n", status)
	err := s.writeResponse(w, req, status, body)
	if err != nil {
		log.Println("Error while writing response: " + err.Error())
	}
	log.Println("ServeHTTP returning")
}

func (s Server) writeResponse(w http.ResponseWriter, req *http.Request, status int, body []byte) error {
	w.WriteHeader(status)
	if len(body) > 0 {
		w.Write(body)
	}
	return nil
}

type BaseResponder struct {
}

func (r BaseResponder) Handle(req *http.Request) (int, []byte) {

	return http.StatusOK, nil
}

func (r BaseResponder) ReadBody(req *http.Request) ([]byte, error) {
	buffer, e := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if e != nil {
		log.Printf("Error reading request body, error=%s", e.Error())
		return nil, e
	}
	return buffer, nil
}

//////////////////////////////////////////
// Diagnostic responders

type MsgResponder struct {
	s string
	BaseResponder
}

func (r MsgResponder) Init(srvr *Server) {
	r.s = MSG
	srvr.Register("GET", "/", r)
	log.Printf("MsgResponder initialized, msg=%s", r.s)
}

func (r MsgResponder) Handle(req *http.Request) (int, []byte) {
	log.Println("MsgResponder.Handle")
	return http.StatusOK, []byte(r.s)
}

type Downer struct {
	server *Server
	BaseResponder
}

func (r Downer) Init(srvr *Server) {
	r.server = srvr
	srvr.Register("GET", BYE_PATH, r)
}

func (r Downer) Handle(req *http.Request) (int, []byte) {
	log.Println("Downer.Handle")
	r.server.Close()
	return http.StatusOK, []byte(BYE_MSG)
}
