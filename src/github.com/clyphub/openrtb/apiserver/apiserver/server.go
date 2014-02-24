/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Lightweight container for Programmatic TV API service
 */
package apiserver

import (
	"encoding/json"
	"github.com/clyphub/openrtb/apiserver/store"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

const MSG = "Hey dere, dis is Hibbing callin\n"
const BYE_PATH = "/bye"
const BYE_MSG = "Buh-bye"

/*
Wrapper for dispatcher/mux
 */
type Router struct {
	routes    []*route
}

type route struct {
	method string
	path string
	handler MethodHandler
}

type Server struct {
	Router
	listener net.Listener
	store *store.Store
}

func NewServer() *Server {
	srvr := new(Server)
	s := new(store.Store)
	srvr.store = s
	return srvr
}

type MethodHandler interface {
	Handle(req *http.Request) (status int, body []byte)
}

type Responder interface {
	MethodHandler
	Init(srvr *Server, store *store.Store)
}

func (s *Server) Init() {
	http.DefaultServeMux.Handle("/", s)
	mr := new(MsgResponder)
	mr.Init(s, s.store)

	sr := new(Downer)
	sr.Init(s, s.store)
}

func (s *Server) Open(laddr string) {
	log.Println("Server opening")
	l, err := net.Listen("tcp", laddr)
	if (err != nil) {
		log.Fatal("Listen: ", err)
	}
	s.listener = l
	http.Serve(s.listener, nil)

}

func (s *Server) Close() {
	if (s.listener != nil) {
		s.listener.Close()
	}
}

func (r *Router) register(method string, path string, handler MethodHandler) {
	route := newRoute(method, path, handler)
	r.routes = append(r.routes, route)
}

func newRoute(method string, path string, handler MethodHandler) *route {
	r := route{method, path, handler}
	return &r
}

func (r *Router) resolveHandler(method string, path string) MethodHandler {
	for _, route := range r.routes {
		ok := route.match(method, path)
		if ok {
			return route.handler
		}
	}
	return nil
}

func (r route) match(method string, path string) bool {
	if(len(method) == 0){
		return false
	}
	if(method != r.method){
		return false
	}
	if(path != r.path){
		return false
	}
	return true
}

func (srvr *Server) AddResponder(responder Responder){
	responder.Init(srvr,  srvr.store)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHTTP entered")
	method := req.Method
	path := req.URL.Path

	status := http.StatusBadRequest
	var body []byte

	sink := r.resolveHandler(method, path)
	if(sink == nil){
		status,body = r.HandleBadMethod(req)
	} else {
		status,body = sink.Handle(req)
	}

	log.Println("Writing response")
	err := writeResponse(w, req, status, body)
	if (err != nil) {
		log.Println("Error while writing response: " + err.Error())
	}
	log.Println("ServeHTTP returning")
}

func writeResponse(w http.ResponseWriter, req *http.Request, status int, body []byte) error {
	w.WriteHeader(status)
	if(len(body)>0){
		w.Write(body)
	}
	return nil
}

type BaseResponder struct {

}

func (r *BaseResponder) Handle(req *http.Request) (status int, body []byte) {

	return http.StatusOK, nil
}


func (r *Router) HandleBadMethod(req *http.Request) (status int, body []byte)  {
	return http.StatusMethodNotAllowed, nil
}

func (r *BaseResponder) readBody(req *http.Request) (b []byte, code int){
	var buffer []byte
	var e error
	buffer, e = ioutil.ReadAll(req.Body)
	if (e != nil) {
		log.Printf("Error reading request body, error=%s",e.Error())
		return nil, http.StatusInternalServerError
	}
	return buffer, http.StatusOK
}

func (r *BaseResponder) marshal(response interface{}) (status int, body []byte){
	// Marshal the results into the response body
	var buffer []byte
	var e error
	buffer, e = json.Marshal(response)
	if (e != nil) {
		log.Printf("Error marshalling response body (%s) error=%s",buffer,e.Error())
		return http.StatusInternalServerError, nil
	}
	return http.StatusOK, buffer
}

type MsgResponder struct {
	s string
	BaseResponder
}

func (r *MsgResponder) Init(srvr *Server, store *store.Store) {
	r.s = MSG
	srvr.register("GET", "/", r)
	log.Printf("MsgResponder initialized, msg=%s", r.s)
}

func (r *MsgResponder) Handle(req *http.Request) (status int, body []byte) {
	log.Println("MsgResponder.Handle")
	return http.StatusOK, []byte(r.s)
}


type Downer struct {
	server *Server
	BaseResponder
}

func (r *Downer) Init(srvr *Server, store *store.Store) {
	r.server = srvr
	srvr.register("GET", BYE_PATH, r)
}

func (r *Downer) Handle(req *http.Request) (status int, body []byte) {
	log.Println("Downer.Handle")
	r.server.Close()
	return http.StatusOK, []byte(BYE_MSG)
}
