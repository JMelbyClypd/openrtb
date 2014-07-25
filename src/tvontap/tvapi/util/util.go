/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: General utilities
*/
package util

import (
	"regexp"
)

func Munge(keys []string) string {
	ret := ""
	lk := len(keys)
	for i, key := range keys {
		ret = ret + key
		if i < lk-1 {
			ret = ret + "/"
		}

	}
	return ret
}

func Unmunge(key string) []string {
	rx := regexp.MustCompile("/")
	arr := rx.Split(key, -1)
	ret := make([]string, 0)
	for _, tok := range arr {
		if len(tok) > 0 {
			ret = append(ret, tok)
		}
	}
	return ret
}
