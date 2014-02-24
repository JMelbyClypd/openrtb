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

type Storable interface {
	GetKey() string
}

type BidRequestObject struct {
	Id                 string             `json:"id"`
	Impressions        []ImpressionObject `json:"imp"`
	Site               SiteObject         `json:",omitempty"`
	App                AppObject          `json:",omitempty"`
	Environment        EnvironmentObject  `json:"env,omitempty"`
	Device             DeviceObject       `json:",omitempty"`
	Audience           AudienceObject     `json:",omitempty"`
	AuctionType        int                `json:"at,omitempty"`
	MaxResponseTime    int                `json:"tmax,omitempty"`
	BuyerSeats         []string           `json:"wseat,omitempty"`
	AllowedCurrencies  []string           `json:"cur,omitempty"`
	BlockedCategories  []string           `json:"bcat,omitempty"`
	BlockedAdvertisers []string           `json:"badv,omitempty"`
	Extension          ExtensionObject    `json:"ext,omitempty"`
}

func(o *BidRequestObject) GetKey() string {
	return o.Id
}

type ImpressionObject struct {
	Id               string          `json:"id"`
	Video            VideoObject     `json:",omitempty"`
	TagId            string          `json:"tagid,omitempty"`
	BidFloor         float32         `json:"bidfloor"`
	BidFloorCurrency string          `json:"bidfloorcur,omitempty"`
	Dateof           time.Time       `json:",omitempty"`
	Extension        ExtensionObject `json:"ext,omitempty"`
}

type SiteObject struct {
}

type AppObject struct {
}

type DeviceObject struct {
}

type VideoObject struct {
	MimeTypes         []string        `json:"mimes"`
	Linearity         int             `json:"linearity"`
	MinDuration       int             `json:"minduration"`
	MaxDuration       int             `json:"maxduration"`
	Protocol          int             `json:"protocol"`
	Sequence          int             `json:"sequence,omitempty"`
	BlockedAttributes []int           `json:"battr,omitempty"`
	Tv                TVObject        `json:"tv,omitempty"`
	Extension         ExtensionObject `json:"ext,omitempty"`
}

type TVObject struct {
	Resolutions []string        `json:"resolutions,omitempty"`
	Aspects     []int           `json:"aspects,omitempty"`
	Broadcast   int             `json:"broadcast,omitempty"`
	Extension   ExtensionObject `json:"ext,omitempty"`
}

type EnvironmentObject struct {
	Id         string          `json:"id,omitempty"`
	Network    string          `json:"network,omitempty"`
	Daypart    string          `json:"daypart,omitempty"`
	DaysOfWeek []string        `json:"daysofweek,omitempty"`
	Categories []string        `json:"cat,omitempty"`
	Producer   ProducerObject  `json:"producer,omitempty"`
	Content    ContentObject   `json:"content,omitempty"`
	Publisher  PublisherObject `json:"publisher,omitempty"`
	Keywords   []string        `json:"keywords,omitempty"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}

type ProducerObject struct {
	Id         string          `json:"id,omitempty"`
	Name       string          `json:"name,omitempty"`
	Categories []string        `json:"cat,omitempty"`
	Domain     url.URL         `json:"domain,omitempty"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}


func(o *ProducerObject) GetKey() string {
	return o.Id
}

type ContentObject struct {
	Id            string          `json:"id,omitempty"`
	Airdate       string          `json:"airdate,omitempty"`
	Episode       int             `json:"episode,omitempty"`
	Title         string          `json:"title,omitempty"`
	Series        string          `json:"series,omitempty"`
	Season        string          `json:"season,omitempty"`
	Url           url.URL         `json:"url,omitempty"`
	Categories    []string        `json:"cat,omitempty"`
	Keywords      []string        `json:"keywords,omitempty"`
	ContentRating string          `json:"contentrating,omitempty"`
	Live          bool            `json:"livestream,omitempty"`
	Length        int             `json:"len,omitempty"`
	Language      string          `json:"language,omitempty"`
	Extension     ExtensionObject `json:"ext,omitempty"`
}


func(o *ContentObject) GetKey() string{
	return o.Id
}

type PublisherObject struct {
	Id        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Domain    url.URL         `json:"domain,omitempty"`
	Extension ExtensionObject `json:"ext,omitempty"`
}


func(o *PublisherObject) GetKey() string{
	return o.Id
}

type AudienceObject struct {
	Count     int             `json:"count,omitempty"`
	Grp       float32         `json:"grp,omitempty"`
	Data      []DataObject    `json:"data,omitempty"`
	Geo       GeoObject       `json:"geo,omitempty"`
	Extension ExtensionObject `json:"ext,omitempty"`
}

type DataObject struct {
	Id        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Segments  []SegmentObject `json:"segment,omitempty"`
	Extension ExtensionObject `json:"ext,omitempty"`
}

type GeoObject struct {
	Latitude      float32         `json:"lat,omitempty"`
	Longitude     float32         `json:"lon,omitempty"`
	Country       string          `json:"country,omitempty"`
	Region        string          `json:"region,omitempty"`
	RegionFIPS104 string          `json:"regionfips104,omitempty"`
	Metro         string          `json:"metro,omitempty"`
	City          string          `json:"city,omitempty"`
	Zipcode       string          `json:"zip,omitempty"`
	Geotype       int             `json:"type,omitempty"`
	Extension     ExtensionObject `json:"ext,omitempty"`
}

type SegmentObject struct {
	Id        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Value     string          `json:"value,omitempty"`
	Ratio     int             `json:"ratio,omitempty"`
	Extension ExtensionObject `json:"ext,omitempty"`
}

type BidResponseObject struct {
	Id          string          `json:"id"`
	SeatBids    []SeatBidObject `json:"seatbid"`
	BidId       string          `json:"bidid,omitempty"`
	BidCurrency string          `json:"cur,omitempty"`
	CustomData  string          `json:"customdata,omitempty"`
	Extension   ExtensionObject `json:"ext,omitempty"`
}

type SeatBidObject struct {
	Bids       []BidObject     `json:"bid,omitempty"`
	Seat       string          `json:"seat,omitempty"`
	WinAsGroup int             `json:"group,omitempty"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}

type BidObject struct {
	Id                 string          `json:"id"`
	ImpressionId       string          `json:"impid"`
	Price              float32         `json:"price:"`
	AdId               string          `json:"adid,omitempty"`
	NotificationUrl    url.URL         `json:"nurl,omitempty"`
	AdMarkup           string          `json:"adm,omitempty"`
	Advertiser         []string        `json:"adomain,omitempty"`
	CampaignId         string          `json:"cid,omitempty"`
	CreativeId         string          `json:"crid,omitempty"`
	CreativeAttributes []string        `json:"attr,omitempty"`
	Extension          ExtensionObject `json:"ext,omitempty"`
}


func(o *BidObject) GetKey() string{
	return o.Id
}

type ExtensionObject struct {
}
