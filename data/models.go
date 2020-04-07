package data

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DataPoint struct {
	gorm.Model
	Province        string    `json:"province"`
	Country         string    `json:"country" gorm:"index:idx_date_country_code"`
	CountryCode     string    `json:"countryCode" gorm:"index:idx_date_country_code"`
	Confirmed       int       `json:"confirmed"`
	Deaths          int       `json:"deaths"`
	Date            time.Time `json:"date" gorm:"index:idx_date_country_code"`
	Lat             string    `json:"lat"`
	Long            string    `json:"long"`
	ExternalCountry string    `json:"-"`
}

var Point DataPoint

type CountyData struct {
	gorm.Model
	FipsId    string    `json:"fipsId" gorm:"index:idx_date_state_county_fips"`
	State     string    `json:"state" gorm:"index:idx_date_state_county_fips"`
	StateCode string    `json:"stateCode"`
	County    string    `json:"county" gorm:"index:idx_date_state_county_fips"`
	Confirmed int       `json:"confirmed"`
	Deaths    int       `json:"deaths"`
	Date      time.Time `json:"date" gorm:"index;index:idx_date_state_county_fips"`
	Lat       string    `json:"lat"`
	Long      string    `json:"long"`
}

var CountyCases CountyData

type CountyHistorical struct {
	gorm.Model
    Data string `sql:"type:text;"`
}

var CountyHist CountyHistorical

type StateHistorical struct {
	gorm.Model
    Data string `sql:"type:text;"`
}

var StateHist StateHistorical

type WorldHistorical struct {
	gorm.Model
    Data string `sql:"type:text;"`
}

var WorldHist WorldHistorical

type ProvinceHistorical struct {
	gorm.Model
    Data string `sql:"type:text;"`
}

var ProvinceHist ProvinceHistorical
