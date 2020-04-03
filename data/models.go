package data

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DataPoint struct {
	gorm.Model
	State           string    `json:"state"`
    Country         string    `json:"country" gorm:"index:idx_date_country_code"`
    CountryCode     string    `json:"countryCode" gorm:"index:idx_date_country_code"`
	Confirmed       int       `json:"confirmed"`
	Deaths          int       `json:"deaths"`
    Date            time.Time `json:"date" gorm:"index:idx_date_country_code"`
	Lat             string    `json:"lat"`
	Long            string    `json:"long"`
	ExternalState   string    `json:"-"`
	ExternalCountry string    `json:"-"`
}

var Point DataPoint

type CountyData struct {
	gorm.Model
	ExternalId string    `json:"externalId"`
	FipsId     string    `json:"fipsId" gorm:"index:idx_date_state_county_fips"`
	State      string    `json:"state" gorm:"index:idx_date_state_county_fips"`
	StateCode  string    `json:"stateCode"`
	County     string    `json:"county" gorm:"index:idx_date_state_county_fips"`
	CountyKey  string    `json:"countyKey"`
	Confirmed  int       `json:"confirmed"`
	Recovered  int       `json:"recovered"`
	Deaths     int       `json:"deaths"`
    Date       time.Time `json:"date" gorm:"index;index:idx_date_state_county_fips"`
}

var CountyCases CountyData
