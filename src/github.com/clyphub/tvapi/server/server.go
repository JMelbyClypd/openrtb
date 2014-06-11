/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Lightweight container for Programmatic TV API service
 */
package server

import (
	"github.com/clyphub/tvapi/store"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
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
	store store.ObjectStore
}

func NewServer() *Server {
	return &Server{store: &store.MapStore{}}
}

type MethodHandler interface {
	Handle(req *http.Request) (status int, body []byte)
}

type Responder interface {
	MethodHandler
	Init(srvr *Server, store store.ObjectStore)
}

func (s *Server) Init() {
	http.DefaultServeMux = http.NewServeMux()
	http.DefaultServeMux.Handle("/", s)
	mr := &MsgResponder{}
	mr.Init(s, s.store)

	sr := &Downer{}
	sr.Init(s, s.store)
}

func (s *Server) Open(laddr string) {
	log.Println("Server opening")
	var err error
	s.listener, err = net.Listen("tcp", laddr)
	if (err != nil) {
		log.Fatal("Listen: ", err)
	}
	http.Serve(s.listener, nil)

}

func (s *Server) Close() {
	if (s.listener != nil) {
		s.listener.Close()
	}
}

func (r *Router) Register(method string, path string, handler MethodHandler) {
	// convert the method to uppercase
	method = strings.ToUpper(method)
	// convert the path to lowercase
	path = strings.ToLower(path)
	log.Printf("Registering method %s to path %s", method, path)
	route := newRoute(method, path, handler)
	r.routes = append(r.routes, route)
}

func newRoute(method string, path string, handler MethodHandler) *route {
	return &route{method, path, handler}
}

func (r *Router) resolveHandler(method string, path string) MethodHandler {
	// convert the method to uppercase
	method = strings.ToUpper(method)
	// convert the path to lowercase
	path = strings.ToLower(path)
	for _, route := range r.routes {
		ok := route.match(method, path)
		if ok {
			return route.handler
		}
	}
	return nil
}

func (r *route) match(method string, path string) bool {
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

	log.Printf("Method=%s",method)
	sink := r.resolveHandler(method, path)
	if(sink == nil){
		status,body = r.HandleBadMethod(req)
	} else {
		status,body = sink.Handle(req)
	}

	log.Println("Writing response")
	err := r.writeResponse(w, req, status, body)
	if (err != nil) {
		log.Println("Error while writing response: " + err.Error())
	}
	log.Println("ServeHTTP returning")
}

func (r *Router) writeResponse(w http.ResponseWriter, req *http.Request, status int, body []byte) error {
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

func (r *BaseResponder) ReadBody(req *http.Request) (b []byte, e error){
	buffer, e := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if (e != nil) {
		log.Printf("Error reading request body, error=%s",e.Error())
		return nil, e
	}
	return buffer, nil
}

func (r *BaseResponder) GetObject() interface{} {
	return nil
}


//////////////////////////////////////////
// Diagnostic responders

type MsgResponder struct {
	s string
	BaseResponder
}

func (r *MsgResponder) Init(srvr *Server, store store.ObjectStore) {
	r.s = MSG
	srvr.Register("GET", "/", r)
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

func (r *Downer) Init(srvr *Server, store store.ObjectStore) {
	r.server = srvr
	srvr.Register("GET", BYE_PATH, r)
}

func (r *Downer) Handle(req *http.Request) (status int, body []byte) {
	log.Println("Downer.Handle")
	r.server.Close()
	return http.StatusOK, []byte(BYE_MSG)
}



