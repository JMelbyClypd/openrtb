/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Test utilities
*/
package testutil

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

/*****************************************
Helper functions
*/

func do(t *testing.T, mthd string, addr string, path string, body string, bodyType string, expectedCode int,
	expected string, checker func(ref string, val string) bool, stopper func()) {
	url := "http://" + addr + path

	req, err := http.NewRequest(mthd, url, strings.NewReader(body))

	if len(bodyType) != 0 {
		req.Header.Add("Content-Type", bodyType)
	}

	if err != nil {
		t.Log("Could not create request")
		t.Fail()
		stopper()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Log("Could not send request")
		t.Fail()
		stopper()
		return
	}

	status := resp.StatusCode
	if status != expectedCode {
		t.Logf("Incorrect response status: %d", status)
		t.Fail()
		stopper()
		return
	}

	var buffer []byte
	buffer, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("Couldn't read body")
		t.Fail()
		stopper()
		return
	}
	if len(buffer) > 0 {
		msg := string(buffer)
		var passes = true
		if checker != nil {
			passes = checker(expected, msg)
		}

		if passes != true {
			t.Fail()
		} else {

		}
	} else {
		t.Log("OK")
	}

	stopper()
}

func DoGet(t *testing.T, addr string, path string, expectedCode int, expected string, checker func(ref string, val string) bool, stopper func()) {
	do(t, "GET", addr, path, "", "", expectedCode, expected, checker, stopper)
}

func DoPost(t *testing.T, addr string, path string, body string, bodyType string, expectedCode int, expected string, checker func(ref string, val string) bool, stopper func()) {
	do(t, "POST", addr, path, body, bodyType, expectedCode, expected, checker, stopper)
}

func DoPut(t *testing.T, addr string, path string, body string, bodyType string, expectedCode int, expected string, checker func(ref string, val string) bool, stopper func()) {
	do(t, "PUT", addr, path, body, bodyType, expectedCode, expected, checker, stopper)
}

func DoDelete(t *testing.T, addr string, path string, expectedCode int, stopper func()) {
	do(t, "DELETE", addr, path, "", "", expectedCode, "", func(ref string, val string) bool { return true }, stopper)
}
