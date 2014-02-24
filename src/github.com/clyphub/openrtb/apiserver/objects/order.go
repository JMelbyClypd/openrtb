/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Data objects for Programmatic TV API test service
 */
package objects

import (
	"net/url"
	"time"
)

type RFPObject struct {
	Id                string          `json:"id"`
	BuyerId           string          `json:"buyerid"`
	Name              string          `json:"name,omitempty"`
	Start             time.Time       `json:"start,omitempty"`
	End               time.Time       `json:"end,omitempty"`
	Advertiser        string          `json:"adomain,omitempty"`
	MaxCPM            float32         `json:"maxcpm,omitempty"`
	MaxSpend          int             `json:"maxspend,omitempty"`
	MinImpressions    int             `json:"minimpressions,omitempty"`
	MaxImpressions    int             `json:"maximpressions,omitempty"`
	Separation        []string        `json:"separation,omitempty"`
	Expiration        time.Time       `json:"expiration,omitempty"`
	DayAndDayparts    []string        `json:"day_dayparts,omitempty"`
	ExcludedProducers []string        `json:"excprods,omitempty"`
	ExcludedGenres    []string        `json:"excgenres,omitempty"`
	Genres            []string        `json:"genres,omitempty"`
	Audience          AudienceObject  `json:"audience"`
	ResponseUrl       url.URL         `json:"purl,omitempty"`
	Extension         ExtensionObject `json:"ext,omitempty"`
}


func(o *RFPObject) GetKey() string{
	return o.Id
}

type ProposalObject struct {
	Id             string          `json:"id"`
	Rfpid          string          `json:"rfpid"`
	BuyerId        string          `json:"buyerid"`
	invalidreasons []string        `json:"invalidreasons,omitempty"`
	Start          time.Time       `json:"start,omitempty"`
	End            time.Time       `json:"end,omitempty"`
	Advertiser     string          `json:"adomain,omitempty"`
	MinCPM         float32         `json:"mincpm,omitempty"`
	MinSpend       int             `json:"minspend,omitempty"`
	MinImpressions int             `json:"minimpressions"`
	MaxImpressions int             `json:"maximpressions"`
	Expiration     time.Time       `json:"expiration,omitempty"`
	Ads            []AdObject      `json:"ads,omitempty"`
	Audience       AudienceObject  `json:"audience"`
	Extension      ExtensionObject `json:"ext,omitempty"`
}


func(o *ProposalObject) GetKey() string{
	return o.Id
}

type OrderObject struct {
	Id         string          `json:"id"`
	RfpId      string          `json:"rfpid"`
	BuyerId    string          `json:"buyerid"`
	Name       string          `json:"name,omitempty"`
	Version    int             `json:"ver"`
	Start      time.Time       `json:"start,omitempty"`
	End        time.Time       `json:"end,omitempty"`
	Advertiser string          `json:"adomain,omitempty"`
	Cpm        float32         `json:"cpm,omitempty"`
	Spend      int             `json:"spend,omitempty"`
	Expiration time.Time       `json:"expiration,omitempty"`
	Ads        []AdObject      `json:"ads"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}


func(o *OrderObject) GetKey() string{
	return o.Id
}

type AdObject struct {
	Id                string          `json:"id"`
	OrderId           string          `json:"orderid"`
	AdId              string          `json:"adid,omitempty"`
	Cpm               float32         `json:"cpm,omitempty"`
	Spend             int             `json:"spend,omitempty"`
	MinImpressions    int             `json:"minimpressions,omitempty"`
	MaxImpressions    int             `json:"maximpressions,omitempty"`
	Separation        []string        `json:"separation,omitempty"`
	DayAndDayparts    []string        `json:"day_dayparts,omitempty"`
	ExcludedProducers []string        `json:"excprods,omitempty"`
	ExcludedGenres    []string        `json:"excgenres,omitempty"`
	Genres            []string        `json:"genres,omitempty"`
	Categories        []string        `json:"cat,omitempty"`
	Audience          AudienceObject  `json:"audience"`
	NotificationUrl   url.URL         `json:"nurl,omitempty"`
	Extension         ExtensionObject `json:"ext,omitempty"`
}


func(o *AdObject) GetKey() string{
	return o.Id
}

type AcceptanceObject struct {
	Id             string          `json:"id"`
	Rfpid          string          `json:"rfpid"`
	BuyerId        string          `json:"buyerid"`
	Accepted       int             `json:"accepted"`
	Invalidreasons []string        `json:"invalidreasons,omitempty"`
	OrderType      int             `json:"type,omitempty"`
	Extension      ExtensionObject `json:"ext,omitempty"`
}

type NotificationObject struct {
	Id           string              `json:"id"`
	OrderId      string              `json:"orderid"`
	BuyerId      string              `json:"buyerid"`
	AdId         string              `json:"adid"`
	Status       int                 `json:"status"`
	Timestamp    time.Time           `json:"timestamp"`
	Measurements []MeasurementObject `json:"measurements,omitempty"`
	OrderType    int                 `json:"type,omitempty"`
	Extension    ExtensionObject     `json:"ext,omitempty"`
}

type MeasurementObject struct {
	Count     int             `json:"count"`
	Window    []string        `json:"window"`
	Geo       GeoObject       `json:"geo,omitempty"`
	Data      []DataObject    `json:"data,omitempty"`
	Extension ExtensionObject `json:"ext,omitempty"`
}
