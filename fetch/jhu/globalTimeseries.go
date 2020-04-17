package jhu

import (
	"covid-tracker/data"
)

type GlobalTimeseries struct{}

var GlobalTS = GlobalTimeseries{}

func (_ *GlobalTimeseries) ToDataPoint(columns []string) *data.DataPoint {
    if columns[1] == "US" { return nil }
    country, province, countryCode := normalizeCountry(columns[1], columns[0])
    lat := columns[2]
    long := columns[3]

    return &data.DataPoint{
        Country: country,
        CountryCode: countryCode,

        Province: province,

        Lat: lat,
        Long: long,

        Population: OverrideForProvince(country, province).Population,
    }
}

func (_ *GlobalTimeseries) TimeColumns(columns []string, _ string) []string {
    return columns[4:]
}
