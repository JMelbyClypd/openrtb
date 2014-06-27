/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Test client for Programmatic TV API service
 */
package client

import (
	"bytes"
	"fmt"
	"github.com/clyphub/tvapi/objects"
	"log"
	"net/http"
	"errors"
	"io/ioutil"
)

const (
	PTV_HDR = "x-ptv-version"
	PTV_HDR_VAL = "1.0"
	CONTENT_TYPE = "application/json"
)

type Client struct {
	http.Client
	baseUrl string
}

func NewClient(srvr string)(*Client, error){
	if(len(srvr) == 0){
		return nil,errors.New("No base URL specified")
	}
	return &Client{baseUrl: srvr}, nil
}

func (c Client) makeUrl(path string) string {
	return "http://" + c.baseUrl + path
}

func (c Client) PostRequest(obj interface{}, path string) error {
	buffer, err := objects.Marshal(obj)
	if (err != nil) {
		log.Printf("Error marshalling request body, error=%s",err.Error())
		return err
	}
	req, err := http.NewRequest("POST", c.makeUrl(path), bytes.NewReader(buffer))
	if(err != nil){
		return err
	}
	req.Header.Set("Content-Type", CONTENT_TYPE)
	req.Header.Set(PTV_HDR, PTV_HDR_VAL)

	resp, err := c.Do(req)
	if(err != nil){
		return err
	}
	status := resp.StatusCode
	if(status == 200){
		return nil
	}
	return fmt.Errorf("Unsuccessful request, server returned %d", status)

}

func (c Client) GetRequest(url string, ref interface{}) error {
	if (len(url) == 0) {
		return errors.New("Could not GET: no URL specified")
	}
	req, err := http.NewRequest("GET", url, nil)
	if(err != nil){
		return err
	}
	req.Header.Set(PTV_HDR, PTV_HDR_VAL)

	resp, err := c.Do(req)
	if(err != nil){
		return err
	}
	status := resp.StatusCode
	if(status == 200){
		buffer, err2 := ioutil.ReadAll(req.Body)
		if(err2 != nil){
			return err2
		}
		err2 = objects.Unmarshal(ref, buffer)
		return err2
	}
	return fmt.Errorf("Unsuccessful request, server returned %d", status)
}

func (c Client) DelRequest(url string) error {
	if (len(url) == 0) {
		return errors.New("Could not GET: no URL specified")
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if(err != nil){
		return err
	}
	req.Header.Set(PTV_HDR, PTV_HDR_VAL)

	resp, err := c.Do(req)
	if(err != nil){
		return err
	}
	status := resp.StatusCode
	if(status == 200){
		return nil
	}
	return fmt.Errorf("Unsuccessful request, server returned %d", status)
}



