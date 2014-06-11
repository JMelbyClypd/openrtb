/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: main package for Programmatic TV API test service examples

 */
package main

import (
	"encoding/json"
	"fmt"
	"github.com/clyphub/tvapi/objects"
	"log"
)

func main() {
	ar := objects.AvailabilityRequestObject{RequestId:"1234"}
	arr := objects.AvailabilityResponseObject{RequestId:"1234"}
	_ = dump(ar)
	_ = dump(arr)
}

func dump(msg interface{}) error {
	buffer, err := toJson(msg)
	if (err != nil) {
		log.Printf("Error marshalling message body (%s) error=%s",buffer,err.Error())
		return err
	}
	s := string(buffer)
	fmt.Println(s)
	return nil
}

func toJson(msg interface{}) (body []byte, e error){
	// Marshal the results into the response body
	buffer, e := json.MarshalIndent(msg,"","    ")
	if (e != nil) {
		log.Printf("Error marshalling response body (%s) error=%s",buffer,e.Error())
		return nil, e
	}
	return buffer, nil
}
