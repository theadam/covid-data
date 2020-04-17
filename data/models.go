package data

import (
	"time"
)

type DataPoint struct {
	Country     string
	CountryCode string

	Province string

	County string
	FipsId string

	Lat  string
	Long string

    key *DataPointKey
    locationKey *LocationKey

	Date       time.Time
	Confirmed  int
	Deaths     int
	Population int
}

type DataPointKey struct {
    Date time.Time
    Country string
    Province string
    County string
}

type LocationKey struct {
    Country string
    Province string
    County string
}

func (dp *DataPoint) Key() *DataPointKey {
    if dp.key != nil {
        return dp.key
    }
    key := &DataPointKey{
        Date: dp.Date,
        Country: dp.Country,
        Province: dp.Province,
        County: dp.County,
    }
    dp.key = key;
    return key
}

func (dp *DataPoint) LocationKey() *LocationKey {
    if dp.locationKey != nil {
        return dp.locationKey
    }
    key := &LocationKey{
        Country: dp.Country,
        Province: dp.Province,
        County: dp.County,
    }
    dp.locationKey = key;
    return key
}
