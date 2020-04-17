package jhu

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"strings"
)

type UsTimeseries struct{}

var UsTS = UsTimeseries{}

func (_ *UsTimeseries) ToDataPoint(columns []string) *data.DataPoint {
    country, province, countryCode := normalizeCountry(columns[7], columns[6])

    fips := padFips(strings.Split(columns[4], ".")[0])
	fipsData, ok := FipsMap[fips]

    county := columns[5]
    stateCode := ""

	if !ok {
        fips = ""
	} else {
        county = fipsData.Name
        stateCode = fipsData.StateCode
        province = utils.StateCodes[stateCode]
    }
	lat := columns[8]
	long := columns[9]

    var population int
    if fipsData.Fips != "" {
        province = utils.StateCodes[fipsData.StateCode]
        county = fipsData.Name
        population = OverrideForFips(fipsData.Fips).Population
    } else {
        population = OverrideForProvince(country, province).Population
    }

    return &data.DataPoint{
        Country: country,
        CountryCode: countryCode,

        Province: province,

        County: county,
        FipsId: fips,

        Lat: lat,
        Long: long,

        Population: population,
    }
}

func (_ *UsTimeseries) TimeColumns(columns []string, kind string) []string {
    if kind == "deaths" {
        return columns[12:]
    }
	return columns[11:]
}
