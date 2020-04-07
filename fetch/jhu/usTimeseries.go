package jhu

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"strings"
)

type UsTimeseries struct{}

var UsTS = UsTimeseries{}

func (_ *UsTimeseries) ExtractExtraData(columns []string) map[string]string {
    fips := padFips(strings.Split(columns[4], ".")[0])
	fipsData, ok := FipsMap[fips]
    county := columns[5]
    state := columns[6]
    stateCode := ""
	if !ok {
        fips = ""
	} else {
        county = fipsData.Name
        stateCode = fipsData.StateCode
        state = utils.StateCodes[stateCode]
    }
	country := columns[7]
	lat := columns[8]
	long := columns[9]
	key := columns[10]
	return map[string]string{
		"fips":    fips,
		"county":  county,
		"state":   state,
		"country": country,
		"lat":     lat,
		"long":    long,
		"key":     key,
	}
}

func toCountyData(timeData TimeValue) data.CountyData {
	fields := timeData.ExtraData

	series := data.CountyData{
		FipsId:    fields["fips"],
		State:     fields["state"],
		StateCode: fields["stateCode"],
		County:    fields["county"],
		Confirmed: timeData.Confirmed,
		Deaths:    timeData.Deaths,
		Date:      timeData.Date,
		Lat:       fields["lat"],
		Long:      fields["long"],
	}
	return series
}

func (_ *UsTimeseries) TimeColumns(columns []string, kind string) []string {
    if kind == "deaths" {
        return columns[12:]
    }
	return columns[11:]
}

func (_ *UsTimeseries) Key(fields map[string]string) string {
	return fields["key"]
}

func (_ *UsTimeseries) Skip(_ map[string]string) bool {
    return false
}
