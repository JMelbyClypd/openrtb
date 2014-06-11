package apiclient

import (
	"github.com/clyphub/tvapi/objects"
	"github.com/clyphub/tvapi/client"
	"log"
	"errors"

)

const (
	AVAIL_URL = "Orders/Availability"
	ORDER_URL = "Orders/Order"
)

type APIClient struct {
	client.Client
}

func (c *APIClient) RequestAvailability(req objects.AvailabilityRequestObject) error {
	return c.PostRequest(req, c.baseUrl + AVAIL_URL)
}

func (c *APIClient) PlaceOrder(req objects.OrderObject) error {
	return c.PostRequest(req, c.baseUrl + ORDER_URL)
}

func (c *APIClient) GetOrders(buyerId string) ([]objects.OrderAcceptanceObject, error){
	if(buyerId == nil){
		return nil, errors.New("Can't get orders: no buyerId specified")
	}
	url := c.baseUrl + ORDER_URL + "/" + buyerId
	var objs = []interface{}
	err := c.GetRequest(url,objs)
	if(err != nil){
		return nil, err
	}
	return objs.([]objects.OrderAcceptanceObject), nil
}

func (c *APIClient) GetOrder(buyerId string, orderId string) (objects.OrderAcceptanceObject, error) {
	if(buyerId == nil){
		return nil, errors.New("Can't get order: no buyerId specified")
	}
	if(orderId == nil){
		return nil, errors.New("Can't get order: no orderId specified")
	}
	url := c.baseUrl + ORDER_URL + "/" + buyerId + "/" + orderId
	var obj = interface{}
	err := c.GetRequest(url,obj)
	if(err != nil){
		return nil, err
	}
	return obj.(objects.OrderAcceptanceObject), nil
}

func (c *APIClient) DeleteOrder(buyerId string, orderId string) error {
	if(buyerId == nil){
		return errors.New("Can't get order: no buyerId specified")
	}
	if(orderId == nil){
		return errors.New("Can't get order: no orderId specified")
	}
	url := c.baseUrl + ORDER_URL + "/" + buyerId + "/" + orderId
	return c.DelRequest(url)
}
