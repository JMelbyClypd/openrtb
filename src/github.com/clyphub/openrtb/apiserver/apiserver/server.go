package apiserver

import (
	"encoding/json"
	"github.com/clyphub/openrtb/apiserver/store"
	"io"
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
}

type Server struct {
	router   *Router
	listener net.Listener
	store *store.Store
}

type Responder interface {
	http.Handler
	Init(srvr *Server, router *Router, store *store.Store)
}

func (s *Server) Init() {
	s.store = new(store.Store)
	s.router = newRouter(s)
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

func (r *Router) register(path string, handler http.Handler) {
	http.DefaultServeMux.Handle(path, handler)
}

func newRouter(srvr *Server) *Router {
	r := new(Router)
	mr := new(MsgResponder)
	mr.Init(srvr, r, srvr.store)

	sr := new(Downer)
	sr.Init(srvr, r, srvr.store)

	return r
}

func (srvr *Server) AddResponder(responder Responder){
	responder.Init(srvr, srvr.router, srvr.store)
}

type MethodHandler interface {
	HandleGet(req *http.Request) (status int, body []byte)
	HandlePut(req *http.Request) (status int, body []byte)
	HandlePost(req *http.Request) (status int, body []byte)
	HandleDelete(req *http.Request) (status int, body []byte)
	HandleBadMethod(req *http.Request) (status int, body []byte)
}

type MethodDispatcher struct {
	sink MethodHandler
}

func (r *MethodDispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("ServeHTTP entered")
	method := req.Method

	status := http.StatusBadRequest
	var body []byte

	switch {
	case method == "POST": {
		status,body = r.sink.HandlePost(req)
	}
	case method == "GET": {
		status,body = r.sink.HandleGet(req)
	}
	case method == "DELETE": {
		status,body = r.sink.HandleDelete(req)
	}
	case method == "PUT": {
		status,body = r.sink.HandleBadMethod(req)
	}
	default:status,body = r.sink.HandleBadMethod(req)
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

func (r *MethodDispatcher) HandleGet(req *http.Request) (status int, body []byte) {

	return http.StatusOK, nil
}

func (r *MethodDispatcher) HandleDelete(req *http.Request) (status int, body []byte) {

	return http.StatusOK, nil
}

func (r *MethodDispatcher) HandleBadMethod(req *http.Request) (status int, body []byte)  {
	return http.StatusMethodNotAllowed, nil
}

func (r *MethodDispatcher) readBody(req *http.Request) (b []byte, code int){
	var buffer []byte
	var e error
	buffer, e = ioutil.ReadAll(req.Body)
	if (e != nil) {
		log.Printf("Error reading request body, error=%s",e.Error())
		return nil, http.StatusInternalServerError
	}
	return buffer, http.StatusOK
}

func (r *MethodDispatcher) marshal(response interface{}) (status int, body []byte){
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
	MethodDispatcher
}

func (r *MsgResponder) Init(srvr *Server, resp *Router, store *store.Store) {
	r.s = MSG
	resp.register("/", r)
	r.sink = r
	log.Printf("MsgResponder initialized, msg=%s", r.s)
}

func (r *MsgResponder) HandleGet(req *http.Request) (status int, body []byte) {
	log.Println("MsgResponder.HandleGet")
	return http.StatusOK, []byte(r.s)
}

func (r *MsgResponder) HandlePost(req *http.Request) (status int, body []byte) {
	log.Println("MsgResponder.HandlePost")
	return http.StatusOK, []byte(r.s)
}

func (r *MsgResponder) HandlePut(req *http.Request) (status int, body []byte) {
	log.Println("MsgResponder.HandlePut")
	return http.StatusOK, []byte(r.s)
}


type Downer struct {
	server *Server
}

func (r *Downer) Init(srvr *Server, resp *Router, store *store.Store) {
	r.server = srvr
	resp.register(BYE_PATH, r)
}

func (r *Downer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.server.Close()
	io.WriteString(w, BYE_MSG)
}
