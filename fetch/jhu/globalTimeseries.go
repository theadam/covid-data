package jhu

import (
	"covid-tracker/data"
)

type GlobalTimeseries struct{}

var GlobalTS = GlobalTimeseries{}

func (_ *GlobalTimeseries) ExtractExtraData(columns []string) map[string]string {
    country, province, countryCode := normalizeCountry(columns[1], columns[0])
    return map[string]string {
        "province": province,
        "country": country,
        "countryCode": countryCode,
        "externalCountry": columns[1],
        "lat": columns[2],
        "long": columns[3],
    }
}

func toDataPoint(timeData TimeValue) data.DataPoint {
    fields := timeData.ExtraData
    series := data.DataPoint{
        Province: fields["province"],
        Country: fields["country"],
        CountryCode: fields["countryCode"],
        Deaths: timeData.Deaths,
        Confirmed: timeData.Confirmed,
        ExternalCountry: fields["externalCountry"],
        Lat:       fields["lat"],
        Long:      fields["long"],
        Date: timeData.Date,
        Population: OverrideForProvince(fields["country"], fields["province"]).Population,
    }
    return series
}

func (_ *GlobalTimeseries) TimeColumns(columns []string, _ string) []string {
    return columns[4:]
}

func (_ *GlobalTimeseries) Key(fields map[string]string) string {
    if fields["province"] != "" {
        return fields["province"] +  ", " + fields["country"]
    }
    return fields["country"]
}

func (_ *GlobalTimeseries) Skip(fields map[string]string) bool {
    return skipGlobal(fields["country"], fields["province"])
}
