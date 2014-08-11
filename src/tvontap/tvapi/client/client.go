/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Test client for Programmatic TV API service
*/
package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"tvontap/tvapi/objects"
)

const (
	PTV_HDR      = "x-ptv-version"
	PTV_HDR_VAL  = "1.0"
	CONTENT_TYPE = "application/json"
)

type Client struct {
	http.Client
	baseUrl string
}

func NewClient(srvr string) (*Client, error) {
	if len(srvr) == 0 {
		return nil, errors.New("No base URL specified")
	}
	return &Client{baseUrl: srvr}, nil
}

func (c Client) makeUrl(path string) string {
	return "http://" + c.baseUrl + path
}

func (c Client) ReadBody(resp *http.Response) ([]byte, error) {
	buffer, e := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if e != nil {
		log.Printf("Error reading request body, error=%s", e.Error())
		return nil, e
	}
	return buffer, nil
}

func (c Client) PostRequest(obj interface{}, path string) *objects.CodedError {
	buffer, err := objects.Marshal(obj)
	if err != nil {
		log.Printf("Error marshalling request body, error=%s", err.Error())
		return objects.FromError(err, 500)
	}
	req, err := http.NewRequest("POST", c.makeUrl(path), bytes.NewReader(buffer))
	if err != nil {
		return objects.FromError(err, 500)
	}
	req.Header.Set("Content-Type", CONTENT_TYPE)
	req.Header.Set(PTV_HDR, PTV_HDR_VAL)

	resp, err := c.Do(req)
	if err != nil {
		return objects.FromError(err, 500)
	}
	status := resp.StatusCode
	if status == 200 {
		return nil
	}
	return objects.NewErrorf("Unsuccessful request, server returned %d", status, status)

}

func (c Client) GetRequest(url string, ref interface{}) *objects.CodedError {
	if len(url) == 0 {
		return objects.NewError("Could not GET: no URL specified", 500)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return objects.FromError(err, 500)
	}
	req.Header.Set(PTV_HDR, PTV_HDR_VAL)

	resp, err := c.Do(req)
	if err != nil {
		return objects.FromError(err, 500)
	}
	status := resp.StatusCode
	log.Printf("GET received status code of %d", status)
	if status == 200 {
		buffer, err2 := c.ReadBody(resp)
		log.Printf("GET received body %s", buffer)
		defer resp.Body.Close()
		if err2 != nil {
			return objects.FromError(err2, 500)
		}
		log.Printf("GET response body read, len=%d", len(buffer))
		err2 = objects.Unmarshal(ref, buffer)
		if err2 != nil {
			return objects.FromError(err2, 500)
		}
		return nil
	}
	return objects.NewErrorf("Unsuccessful request, server returned %d", status, status)
}

func (c Client) DelRequest(url string) *objects.CodedError {
	if len(url) == 0 {
		return objects.NewError("Could not GET: no URL specified", 500)
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return objects.FromError(err, 500)
	}
	req.Header.Set(PTV_HDR, PTV_HDR_VAL)

	resp, err := c.Do(req)
	if err != nil {
		return objects.FromError(err, 500)
	}
	status := resp.StatusCode
	if status == 200 {
		return nil
	}
	return objects.NewErrorf("Unsuccessful request, server returned %d", status, status)
}
