package jhu

type UsTerritoryTimeseries struct{}

var UsTerritoryTS = UsTerritoryTimeseries{}

func (_ *UsTerritoryTimeseries) ExtractExtraData(columns []string) map[string]string {
    if (columns[3] == "840") {
        return map[string]string{
            "countryCode": "840",
        }
    }
    country, province, countryCode := normalizeCountry(columns[6], "")
	return map[string]string{
        "province": province,
        "country": country,
        "countryCode": countryCode,
        "externalCountry": columns[6],
        "lat": columns[8],
        "long": columns[9],
	}
}

func (_ *UsTerritoryTimeseries) TimeColumns(columns []string, kind string) []string {
    if kind == "deaths" {
        return columns[12:]
    }
	return columns[11:]
}

func (_ *UsTerritoryTimeseries) Key(fields map[string]string) string {
    if fields["province"] != "" {
        return fields["province"] +  ", " + fields["country"]
    }
    return fields["country"]
}

func (_ *UsTerritoryTimeseries) Skip(fields map[string]string) bool {
    return fields["countryCode"] == "840"
}
