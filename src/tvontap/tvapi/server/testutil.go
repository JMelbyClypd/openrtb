/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Test utilities
*/
package server

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
	expected string, checker func(ref string, val string) bool) {
	url := "http://" + addr + path

	mthd = strings.ToUpper(mthd)
	req, err := http.NewRequest(mthd, url, strings.NewReader(body))

	if len(bodyType) != 0 {
		req.Header.Add("Content-Type", bodyType)
	}

	if err != nil {
		t.Log("Could not create request")
		t.Fail()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Logf("Could not send request to %s: %s", url, err.Error())
		t.Fail()
		return
	}

	status := resp.StatusCode
	if status != expectedCode {
		t.Logf("Incorrect response status: %d", status)
		t.Fail()
		return
	}

	var buffer []byte
	buffer, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Log("Couldn't read body")
		t.Fail()
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

}

func DoGet(t *testing.T, addr string, path string, expectedCode int, expected string, checker func(ref string, val string) bool) {
	do(t, "GET", addr, path, "", "", expectedCode, expected, checker)
}

func DoPost(t *testing.T, addr string, path string, body string, bodyType string, expectedCode int, expected string, checker func(ref string, val string) bool) {
	do(t, "POST", addr, path, body, bodyType, expectedCode, expected, checker)
}

func DoPut(t *testing.T, addr string, path string, body string, bodyType string, expectedCode int, expected string, checker func(ref string, val string) bool) {
	do(t, "PUT", addr, path, body, bodyType, expectedCode, expected, checker)
}

func DoDelete(t *testing.T, addr string, path string, expectedCode int) {
	do(t, "DELETE", addr, path, "", "", expectedCode, "", func(ref string, val string) bool { return true })
}
