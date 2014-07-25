/*
Copyright 2014 clypd, inc.  All rights reserved.

Author: J. Melby

Description: Data objects for Programmatic TV API test service
*/
package objects

import (
	"encoding/json"
	"fmt"
	"log"
)

type CodedError struct {
	msg    string
	status int
}

func (e CodedError) Error() string {
	return e.msg
}

func (e CodedError) Code() int {
	return e.status
}

func NewError(m string, s int) *CodedError {
	return &CodedError{msg: m, status: s}
}

func FromError(e error, s int) *CodedError {
	return &CodedError{msg: e.Error(), status: s}
}

func NewErrorf(format string, status int, a ...interface{}) *CodedError {
	return NewError(fmt.Sprintf(format, a...), status)
}


type Storable interface {
	GetKey() string
}

type Objectable interface {
	GetObject() interface{}
}

func Unmarshal(ref interface{}, buffer []byte) error {
	if ref == nil {
		log.Println("Attempting to unmarshal into an empty pointer")
	}
	log.Printf("Unmarshalling json\n%s\n into object of type %T", buffer, ref)
	e := json.Unmarshal(buffer, ref)
	log.Printf("Unmarshalled object with type %T", ref)
	if e != nil {
		log.Printf("Error unmarshalling, error=%s", e.Error())
		return e
	}
	return nil
}

func Marshal(obj interface{}) ([]byte, error) {
	// Marshal the results into the response body
	buffer, e := json.MarshalIndent(obj, "", "    ")
	if e != nil {
		log.Printf("Error marshalling object body (%s) error=%s", buffer, e.Error())
		return nil, e
	}
	return buffer, nil
}

type ExtensionObject struct {
}

const (
	Sunday    = "Su"
	Monday    = "Mo"
	Tuesday   = "Tu"
	Wednesday = "We"
	Thursday  = "Th"
	Friday    = "Fr"
	Saturday  = "Sa"
)

type EnvironmentObject struct {
	Id         string          `json:"id,omitempty"`
	Network    string          `json:"network,omitempty"`
	Genre      string          `json:"genre,omitempty"`
	Daypart    string          `json:"daypart,omitempty"`
	DaysOfWeek []string        `json:"daysOfWeek,omitempty"`
	Categories []string        `json:"netCategories,omitempty"`
	Producer   string          `json:"producer,omitempty"`
	Content    ContentObject   `json:"content,omitempty"`
	Publisher  string          `json:"publisher,omitempty"`
	TVFormat   TVObject        `json:"tv,omitempty"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}

func (o EnvironmentObject) GetKey() string {
	return o.Id
}

type ContentObject struct {
	Id         string          `json:"id,omitempty"`
	Airdate    string          `json:"airdate,omitempty"`
	Episode    int             `json:"episode,omitempty"`
	Title      string          `json:"title,omitempty"`
	Series     string          `json:"series,omitempty"`
	Season     string          `json:"season,omitempty"`
	Categories []string        `json:"categories,omitempty"`
	TVFormat   TVObject        `json:"tv,omitempty"`
	Keywords   []string        `json:"keywords,omitempty"`
	Ratings    []string        `json:"ratings,omitempty"`
	Length     int             `json:"len,omitempty"`
	Languages  []string        `json:"languages,omitempty"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}

func (o ContentObject) GetKey() string {
	return o.Id
}

type TVObject struct {
	Resolutions []string        `json:"resolutions,omitempty"`
	Aspects     []int           `json:"aspects,omitempty"`
	IsBroadcast bool            `json:"isBroadcast,omitempty"`
	Extension   ExtensionObject `json:"ext,omitempty"`
}

type AudienceObject struct {
	Count      int             `json:"count"`
	Grps       float32         `json:"grps,omitempty"`
	Confidence int             `json:"confidence,omitempty"`
	Data       []DataObject    `json:"data"`
	Geo        GeoObject       `json:"geo,omitempty"`
	Extension  ExtensionObject `json:"ext,omitempty"`
}

type DataObject struct {
	Id         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Attributes []AttributeObject `json:"attributes"`
	Extension  ExtensionObject   `json:"ext,omitempty"`
}

type GeoObject struct {
	Area               []PointObject   `json:"area,omitempty"`
	Country            string          `json:"country,omitempty"`
	Regions            []string        `json:"regions,omitempty"`
	Metros             []string        `json:"metros,omitempty"`
	PoliticalDistricts []string        `json:"politicalDistricts,omitempty"`
	Cities             []string        `json:"cities,omitempty"`
	PostalCodes        []string        `json:"postalCodes,omitempty"`
	SysCodes           []string        `json:"sysCodes,omitempty"`
	Extension          ExtensionObject `json:"ext,omitempty"`
}

type PointObject struct {
	Latitude  float32         `json:"lat,omitempty"`
	Longitude float32         `json:"lon,omitempty"`
	Extension ExtensionObject `json:"ext,omitempty"`
}

type AttributeObject struct {
	Id           string               `json:"id,omitempty"`
	Name         string               `json:"name"`
	Bins         []AttributeBinObject `json:"bins"`
	HouseholdIDs []SpecificHHsObject  `json:"householdIDs,omitempty"`
	Extension    ExtensionObject      `json:"ext,omitempty"`
}

type AttributeBinObject struct {
	Range     string          `json:"range"`
	Index     int             `json:"index"`
	Extension ExtensionObject `json:"ext,omitempty"`
}

type SpecificHHsObject struct {
	MatchProvider string          `json:"matchProvider"`
	IDs           []string        `json:"ids"`
	Extension     ExtensionObject `json:"ext,omitempty"`
}

type AckObject struct {
	OrderId     string          `json:"orderId"`
	BuyerId     string          `json:"buyerId"`
	Version     int             `json:"version"`
	IsCompleted bool            `json:"isCompleted`
	Errors      []string        `json:"errors"`
	Operation   int             `json:"operation"`
	Extension   ExtensionObject `json:"ext,omitempty"`
}

func (o AckObject) GetKey() string {
	return o.OrderId
}
