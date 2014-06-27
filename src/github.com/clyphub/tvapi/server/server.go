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

const (
	MSG = "Hey dere, dis is Hibbing callin\n"
	BYE_PATH = "/bye"
	BYE_MSG = "Buh-bye"
)

type CodedError struct {
	msg string
	status int
}

func (e CodedError) Error() string {
	return e.msg
}

func (e CodedError) Code() int {
	return e.status
}

func NewError(m string, s int) *CodedError {
	return &CodedError{msg: m, status: s}
}

/*
Wrapper for dispatcher/mux
 */
type Router struct {
	routes    map[string]route
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
	s:= Server{store: store.NewMapStore() }
	s.routes = 	make(map[string]route)
	return &s
}

type MethodHandler interface {
	Handle(req *http.Request) (status int, body []byte)
}

type Responder interface {
	MethodHandler
	Init(srvr *Server, store store.ObjectStore)
}

func (s Server) Init() {
	http.DefaultServeMux = http.NewServeMux()
	http.DefaultServeMux.Handle("/", s)

	mr := &MsgResponder{}
	mr.Init(&s, s.store)

	sr := &Downer{}
	sr.Init(&s, s.store)
}

func (s Server) Open(laddr string) {
	log.Println("Server opening")
	var err error
	s.listener, err = net.Listen("tcp", laddr)
	if (err != nil) {
		log.Fatal("Listen: ", err)
	}
	http.Serve(s.listener, nil)

}

func (s Server) Close() {
	if (s.listener != nil) {
		s.listener.Close()
	}
}

func (r Router) toKey(method string, path string) string {
	if(len(method)==0 || len(path)==0){
		log.Println("Attempting to register empty method or path")
		return ""
	}
	return method + ":" + path
}

func (r Router) Register(method string, path string, handler MethodHandler) {

	if(handler == nil){
		log.Printf("Attempted to register nil handler with method %s and path %s", method, path)
		return
	}
	// convert the method to uppercase
	method = strings.ToUpper(method)
	// convert the path to lowercase
	path = strings.ToLower(path)
	key := r.toKey(method, path)
	route := newRoute(method, path, handler)
	r.routes[key] = route
	log.Printf("%d routes", len(r.routes))
}

func newRoute(method string, path string, handler MethodHandler) route {
	return route{method, path, handler}
}

func (r Router) resolveHandler(method string, path string) MethodHandler {
	// convert the method to uppercase
	method = strings.ToUpper(method)
	// convert the path to lowercase
	path = strings.ToLower(path)
	key := r.toKey(method, path)
	route, ok := r.routes[key]
	if(ok){
		log.Printf("Matched route with method %s and path %s", method, path)
		return route.handler
	}
	log.Printf("No handler found for method %s and path %s", method, path)
	return nil
}

func (r Router) match(ro route, method string, path string) bool {
	log.Printf("Matching %s %s", method, path)
	if(len(method) == 0){
		return false
	}
	if(method != ro.method){
		return false
	}
	if(path != ro.path){
		return false
	}
	return true
}

func (srvr Server) AddResponder(responder Responder){
	responder.Init(&srvr,  srvr.store)
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	log.Printf("Router.ServerHTTP entered with method %s and path %s", method, path)

	status := http.StatusBadRequest
	var body []byte

	sink := r.resolveHandler(method, path)
	if(sink == nil){
		status,body = r.HandleBadMethod(req)
	} else {
		status,body = sink.Handle(req)
	}

	log.Printf("ServeHTTP: Writing response with status %d\n", status)
	err := r.writeResponse(w, req, status, body)
	if (err != nil) {
		log.Println("Error while writing response: " + err.Error())
	}
	log.Println("ServeHTTP returning")
}

func (r Router) writeResponse(w http.ResponseWriter, req *http.Request, status int, body []byte) error {
	w.WriteHeader(status)
	if(len(body)>0){
		w.Write(body)
	}
	return nil
}

type BaseResponder struct {

}

func (r BaseResponder) Handle(req *http.Request) (int, []byte) {

	return http.StatusOK, nil
}


func (r Router) HandleBadMethod(req *http.Request) (int, []byte)  {
	return http.StatusMethodNotAllowed, nil
}

func (r BaseResponder) ReadBody(req *http.Request) ([]byte, error){
	buffer, e := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if (e != nil) {
		log.Printf("Error reading request body, error=%s",e.Error())
		return nil, e
	}
	return buffer, nil
}

func (r BaseResponder) GetObject() interface{} {
	return nil
}


//////////////////////////////////////////
// Diagnostic responders

type MsgResponder struct {
	s string
	BaseResponder
}

func (r MsgResponder) Init(srvr *Server, store store.ObjectStore) {
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

func (r Downer) Init(srvr *Server, store store.ObjectStore) {
	r.server = srvr
	srvr.Register("GET", BYE_PATH, r)
}

func (r Downer) Handle(req *http.Request) (int, []byte) {
	log.Println("Downer.Handle")
	r.server.Close()
	return http.StatusOK, []byte(BYE_MSG)
}



