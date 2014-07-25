/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Implementation for Programmatic TV Order API
*/
package impl

import (
	"log"
	"tvontap/tvapi/objects"
	"tvontap/tvapi/server"
)

type GenericGetProcessor struct {
	StoreManager
}

func (ggp GenericGetProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	return nil, nil
}

func (ggp GenericGetProcessor) ValidateRequest(pathTokens []string, queryTokens []string,
	msg objects.Storable) *server.CodedError {
	// Check the path

	// Check the query tokens

	return nil
}

func (ggp GenericGetProcessor) ProcessRequest(pathTokens []string, queryTokens []string, req objects.Storable,
	responder *APIResponder) ([]objects.Storable, *server.CodedError) {

	lenp := len(pathTokens)
	log.Printf("Path parsed to %d tokens, ggp path has %d tokens", lenp, len(ggp.pathKeys))
	switch lenp - len(ggp.pathKeys) {
	case 0:
	{ // Addressed to general resource
		log.Println("GGP Get: general resource")
		return nil, server.NewError("Not Found", 404)
	}
	case 1:
	{ // Addressed to buyerId
		log.Printf("GGP Get: all buyer %s orders", pathTokens[lenp-1])
		ret, err := ggp.GetAllByBuyer(pathTokens[lenp-1])
		if err != nil {
			return nil, server.NewError(err.Error(), 500)
		}
		return ret, nil
	}
	case 2:
	{ // Addressed to specific item
		log.Printf("GGP Get: specific order %s from buyer %s", pathTokens[lenp-1], pathTokens[lenp-2])
		ret, err := ggp.GetObject(pathTokens[lenp-2], pathTokens[lenp-1])
		if err != nil {
			return nil, server.NewError(err.Error(), 500)
		}
		if ret == nil {
			return nil, server.NewError("Not Found", 404)
		}
		retArr := make([]objects.Storable, 1)
		retArr[0] = ret
		return retArr, nil
	}
	}

	return nil, server.NewError("Problem processing GET request - too many path tokens", 400)
}

type GenericDeleteProcessor struct {
	StoreManager
}

func (gdp GenericDeleteProcessor) Unmarshal(buffer []byte) (objects.Storable, error) {
	return nil, nil
}

func (gdp GenericDeleteProcessor) ValidateRequest(pathTokens []string, queryTokens []string,
	msg objects.Storable) *server.CodedError {
	// Check the path

	// Check the query tokens

	return nil
}

func (gdp GenericDeleteProcessor) ProcessRequest(pathTokens []string, queryTokens []string, req objects.Storable,
	responder *APIResponder) ([]objects.Storable, *server.CodedError) {

	lenp := len(pathTokens)
	log.Printf("Path parsed to %d tokens, ggp path has %d tokens", lenp, len(ggp.pathKeys))
	switch lenp - len(gdp.pathKeys) {
	case 0:
	{ // Addressed to general resource
		log.Println("GDP DELETE: general resource")
		return nil, server.NewError("Not Found", 404)
	}
	case 1:
	{ // Addressed to buyerId
		log.Printf("GDP Delete: all buyer %s orders", pathTokens[lenp-1])
		err := gdp.DeleteAllByBuyer(pathTokens[lenp-1])
		if err != nil {
			return nil, server.NewError(err.Error(), 500)
		}
		return nil, nil
	}
	case 2:
	{ // Addressed to specific item
		log.Printf("GDP Delete: specific order %s from buyer %s", pathTokens[lenp-1], pathTokens[lenp-2])
		err := gdp.DeleteObject(pathTokens[lenp-2], pathTokens[lenp-1])
		if err != nil {
			return nil, server.NewError(err.Error(), 500)
		}

		retArr := make([]objects.Storable, 0)
		return retArr, nil
	}
	}

	return nil, server.NewError("Problem processing DELETE request - too many path tokens", 400)
}
