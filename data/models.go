package data

import (
    "github.com/jinzhu/gorm"
	"time"
)

type DataPoint struct {
    gorm.Model
    State string `json:"state"`
    Country string `json:"country"`
    CountryCode string `json:"countryCode"`
    Confirmed int `json:"confirmed"`
    Deaths int `json:"deaths"`
    Date time.Time `json:"date"`
    Lat string `json:"lat"`
    Long string `json:"long"`
    ExternalState string `json:"-"`
    ExternalCountry string `json:"-"`
}

var Point DataPoint

type CountyData struct {
    gorm.Model
    ExternalId string `json:"externalId"`
    State string `json:"state"`
    StateCode string `json:"stateCode"`
    County string `json:"county"`
    Confirmed int `json:"confirmed"`
    Deaths int `json:"deaths"`
    Date time.Time `json:"date"`
}

var CountyCases CountyData
