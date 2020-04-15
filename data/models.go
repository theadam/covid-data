package data

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DataPoint struct {
	gorm.Model
	Province        string    `json:"province"`
	Country         string    `json:"country"`
	CountryCode     string    `json:"countryCode"`
	Confirmed       int       `json:"confirmed"`
	Deaths          int       `json:"deaths"`
	Date            time.Time `json:"date"`
	Lat             string    `json:"lat"`
	Long            string    `json:"long"`
	Population            int    `json:"population"`
	ExternalCountry string    `json:"-"`
}

var Point DataPoint

type CountyData struct {
	gorm.Model
	FipsId    string    `json:"fipsId"`
	State     string    `json:"state"`
	StateCode string    `json:"stateCode"`
	County    string    `json:"county"`
	Confirmed int       `json:"confirmed"`
	Deaths    int       `json:"deaths"`
	Date      time.Time `json:"date"`
	Lat       string    `json:"lat"`
	Population            int    `json:"population"`
	Long      string    `json:"long"`
}

var CountyCases CountyData

