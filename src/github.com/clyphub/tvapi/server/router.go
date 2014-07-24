/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Resource router for Programmatic TV API service
 */
package server

import (
	"errors"
	"log"
	"net/http"
	"github.com/clyphub/tvapi/util"
	"strings"
)

type Node struct {
	getter MethodHandler
	putter MethodHandler
	poster MethodHandler
	deleter MethodHandler
	path string
	children []*Node
}

func newNode(pathToken string) *Node {
	return &Node{path: pathToken, children: make([]*Node, 0,10)}
}

func(n *Node) Add(child *Node)(){
	if(&child == nil){
		return
	}
	n.children = append(n.children,child)
	log.Printf("Node %s added child %s, children=%d", n.path, child.path, len(n.children))
}

func(n *Node) FindChild(token string) *Node {
	for _, ch := range n.children {
		if(ch.path == token) {
			return ch
		}
	}
	return nil
}

func(n *Node) AddPath(path string) (*Node, error) {
	// Need to add code to strip any trailing /

	sp := util.Unmunge(path)
	cnt := len(sp)
	currNode := n
	for i:= 1; i<cnt;i++ {
		nn := currNode.FindChild(sp[i])
		if(nn == nil ){
			log.Printf("Adding node %s to parent %s", sp[i], currNode.path)
			nn = newNode(sp[i])
			currNode.Add(nn)
		}
		currNode = nn
	}
	return currNode, nil
}

func(n *Node) FindLeaf(path string) *Node {
	sp := util.Unmunge(path)
	cnt := len(sp)
	currNode := n
	for i:= 1; i<cnt;i++ {
		nn := currNode.FindChild(sp[i])
		if(nn == nil){
			return currNode
		}
		currNode = nn
	}
	return currNode
}

func(n *Node) SetHandler(method string, handler MethodHandler) error {
	switch(method){
	case "GET": {
		n.getter = handler
		return nil
	}
	case "PUT":{
		n.putter = handler
		return nil
	}
	case "POST":{
		n.poster = handler
		return nil
	}
	case "DELETE":{
		n.deleter = handler
		return nil
	}
	}
	return errors.New("Attempted to set handler for unsupported method " + method)
}

func(n *Node) GetHandler(method string) MethodHandler {
	switch(method){
	case "GET": {
		return n.getter
	}
	case "PUT":{
		return n.putter
	}
	case "POST":{
		return n.poster
	}
	case "DELETE":{
		return n.deleter
	}
	}
	return nil
}

type Router struct {
	root *Node
	badMethodHandler BadMethodHandler
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
	n, e := r.root.AddPath(path)
	if(e != nil){
		log.Println(e.Error())
		return
	}
	e = n.SetHandler(method, handler)
	if(e != nil){
		log.Println(e.Error())
		return
	}
	log.Printf("Registered handler for method %s and path %s", method,path)
}


func (r Router) resolveHandler(method string, path string) MethodHandler {
	// convert the method to uppercase
	method = strings.ToUpper(method)
	// convert the path to lowercase
	path = strings.ToLower(path)

	n := r.root.FindLeaf(path)
	if(n != nil){
		log.Printf("Found node for path %s at %s", path, n.path)
		h := n.GetHandler(method)
		if(h != nil){
			log.Printf("Matched route with method %s and path %s", method, path)
			return h
		}
	}
	log.Printf("No handler found for method %s and path %s", method, path)
	return r.badMethodHandler
}

func(r Router) dumpNodes() {
	level := 0
	r.root.dump(level)
}

func(n *Node) dump(level int) {
	log.Printf("Node: %s", n.path)
	level = level +1
	for _,child := range n.children{
		child.dump(level)
	}
}


type BadMethodHandler struct {

}

func (h BadMethodHandler) Handle(req *http.Request) (status int, body []byte) {
	return http.StatusMethodNotAllowed, nil
}
