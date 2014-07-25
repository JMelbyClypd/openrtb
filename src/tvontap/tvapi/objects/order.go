/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Data objects for Programmatic TV API test service
*/
package objects

import (
	"time"
)

type AvailabilityRequestObject struct {
	RequestId           string          `json:"requestId"`
	BuyerId             string          `json:"buyerId"`
	AdvertiserId        string          `json:"advertiserId,omitempty"`
	Start               time.Time       `json:"start,omitempty"`
	End                 time.Time       `json:"end,omitempty"`
	MaxCPM              float32         `json:"maxCpm,omitempty"`
	MaxSpend            int             `json:"maxSpend,omitempty"`
	MinImpressions      int             `json:"minImpressions"`
	MaxImpressions      int             `json:"maxImpressions,omitempty"`
	Separation          []string        `json:"separation,omitempty"`
	Expiration          time.Time       `json:"expiration,omitempty"`
	Networks            []string        `json:"networks,omitempty"`
	DayAndDayparts      []string        `json:"dayDayparts,omitempty"`
	ExcludedDayDayparts []string        `json:"excludedDayDayparts,omitempty"`
	ExcludedNetworks    []string        `json:"excludedNets,omitempty"`
	ExcludedGenres      []string        `json:"excludedGenres,omitempty"`
	Genres              []string        `json:"genres,omitempty"`
	Audience            AudienceObject  `json:"audience"`
	ResponseUrl         string          `json:"responseUrl,omitempty"`
	Extension           ExtensionObject `json:"ext,omitempty"`
}

func (o AvailabilityRequestObject) GetKey() string {
	return o.RequestId
}

type AvailabilityResponseObject struct {
	RequestId      string            `json:"requestId"`
	BuyerId        string            `json:"buyerId"`
	InvalidReasons []string          `json:"invalidReasons,omitempty"`
	Start          time.Time         `json:"start,omitempty"`
	End            time.Time         `json:"end,omitempty"`
	MinCPM         float32           `json:"minCpm,omitempty"`
	MinSpend       int               `json:"minspend,omitempty"`
	MinImpressions int               `json:"minImpressions"`
	MaxImpressions int               `json:"maxImpressions"`
	Expiration     time.Time         `json:"expiration,omitempty"`
	Placements     []PlacementObject `json:"placements,omitempty"`
	Extension      ExtensionObject   `json:"ext,omitempty"`
}

type PlacementObject struct {
	Id          string              `json:"id"`
	SupplierId  string              `json:"supplierId,omitempty"`
	Environment []EnvironmentObject `json:"environment"`
	Audience    AudienceObject      `json:"audience"`
	Cpm         float32             `json:"cpm,omitempty"`
	UnitPrice   float32             `json:"unitPrice,omitempty"`
	Extension   ExtensionObject     `json:"ext,omitempty"`
}

func (o AvailabilityResponseObject) GetKey() string {
	return o.RequestId
}

type OrderObject struct {
	RequestId    string            `json:"requestId"`
	BuyerId      string            `json:"buyerId"`
	AdvertiserId string            `json:"advertiserId"`
	Name         string            `json:"name,omitempty"`
	Version      int               `json:"version"`
	Start        time.Time         `json:"start,omitempty"`
	End          time.Time         `json:"end,omitempty"`
	Cpm          float32           `json:"cpm,omitempty"`
	TotalCost    int               `json:"totalCost,omitempty"`
	DealId       string            `json:"dealId,omitempty"`
	Expiration   time.Time         `json:"expiration,omitempty"`
	OrderLines   []OrderlineObject `json:"orderLines"`
	Extension    ExtensionObject   `json:"ext,omitempty"`
}

func (o OrderObject) GetKey() string {
	return o.RequestId
}

type OrderlineObject struct {
	Id                  string            `json:"id"`
	Version             int               `json:"version"`
	CreativeId          string            `json:"creativeId,omitempty"`
	Cpm                 float32           `json:"cpm,omitempty"`
	LineCost            int               `json:"lineCost,omitempty"`
	MinImpressions      int               `json:"minImpressions,omitempty"`
	MaxImpressions      int               `json:"maxImpressions,omitempty"`
	MaxTgtImpressions   int               `json:"maxTgtImpressions,omitempty"`
	OverrunPolicy       int               `json:"overrunPolicy,omitempty"`
	Separation          []string          `json:"separation,omitempty"`
	IsCancellable       bool              `json:"isCancellable"`
	IsGuaranteed        bool              `json:"isGuaranteed"`
	DayAndDayparts      []string          `json:"dayDayparts,omitempty"`
	ExcludedDayDayparts []string          `json:"excludedDayDayparts,omitempty"`
	Networks            []string          `json:"networks,omitempty"`
	ExcludedNets        []string          `json:"excludedNets,omitempty"`
	ExcludedGenres      []string          `json:"excludedGenres,omitempty"`
	Genres              []string          `json:"genres,omitempty"`
	AdCategories        []string          `json:"adCategories,omitempty"`
	Environment         EnvironmentObject `json:"environment,omitempty"`
	Audience            AudienceObject    `json:"audience,omitempty"`
	CurrencyPolicy      int               `json:"currPolicy"`
	Currency            []string          `json:"currency,omitempty"`
	ReportUrl           string            `json:"reportUrl,omitempty"`
	Extension           ExtensionObject   `json:"ext,omitempty"`
}

func (o OrderlineObject) GetKey() string {
	return o.Id
}

type OrderAcceptanceObject struct {
	RequestId      string          `json:"requestId"`
	BuyerId        string          `json:"buyerId"`
	Version        int             `json:"version"`
	Accepted       bool            `json:"accepted"`
	Invalidreasons []string        `json:"invalidReasons,omitempty"`
	OrderType      int             `json:"type,omitempty"`
	Cpm            float32         `json:"cpm"`
	TotalCost      int             `json:"totalCost,omitempty"`
	DealId         string          `json:"dealId,omitempty"`
	ECP            int             `json:"ecp,omitmepty"`
	OrderLineIds   []string        `json:"orderLineIds"`
	Extension      ExtensionObject `json:"ext,omitempty"`
}

func (o OrderAcceptanceObject) GetKey() string {
	return o.RequestId
}

type NotificationObject struct {
	OrderId      string              `json:"orderid"`
	BuyerId      string              `json:"buyerid"`
	OrderlineId  string              `json:"orderLineId"`
	CreativeId   string              `json:"creativeId"`
	Status       int                 `json:"status"`
	Timestamp    time.Time           `json:"timestamp"`
	Rollup       MeasurementObject   `json:"rollup,omitempty"`
	Measurements []MeasurementObject `json:"measurements,omitempty"`
	Extension    ExtensionObject     `json:"ext,omitempty"`
}

func (o NotificationObject) GetKey() string {
	return o.OrderId
}

type MeasurementObject struct {
	TotalImpressions  int               `json:"totalImpressions"`
	TargetImpressions int               `json:"targetImpressions"`
	Grps              float32           `json:"grps,omitempty"`
	Source            string            `json:"source"`
	Validators        []string          `json:"validators,omitempty"`
	Window            []string          `json:"window"`
	Placements        []PlacementObject `json:"placements,omitempty"`
	Geo               GeoObject         `json:"geo,omitempty"`
	CumeAudience      AudienceObject    `json:"cumeAudience,omitempty"`
	Extension         ExtensionObject   `json:"ext,omitempty"`
}
