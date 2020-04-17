package jhu

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type OverrideData struct {
	UID         string
	Iso2        string
	Iso3        string
	Code3       string
	Fips        string
	County      string
	CountryCode string
	Province    string
	Country    string
	Lat         string
	Lng         string
	CombinedKey string
	Population  int
}

func fetchOverrides() []OverrideData {
	body := fetchFromRepo("master", "csse_covid_19_data/UID_ISO_FIPS_LookUp_Table.csv")
	defer body.Close()

	reader := csv.NewReader(body)
	data := make([]OverrideData, 0)
    // Skip Header
    reader.Read()
	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}

        country, province, countryCode := normalizeCountry(columns[7], columns[6])
        population, err := strconv.Atoi(columns[11])
		if err != nil {
            if country != "Cruise" && columns[5] != "Unassigned" && !strings.HasPrefix(columns[5], "Out of") && province != "Recovered" {
                fmt.Println("No population for " + columns[0] + ", " + columns[5] + ", " + province + ", " + country)
            }
            population = 0
		}
		data = append(data, OverrideData{
			UID:         columns[0],
			Iso2:        columns[1],
			Iso3:        columns[2],
			Code3:       columns[3],
			Fips:        padFips(columns[4]),
			County:      columns[5],
			Province:    province,
            Country: country,
            CountryCode: countryCode,
			Lat:         columns[8],
			Lng:         columns[9],
			CombinedKey: columns[10],
			Population:  population,
		})
	}
    return data
}
var overridedata = fetchOverrides()

func OverrideFromUID(uid string) OverrideData {
    for _, item := range(overridedata) {
        if item.UID == uid {
            return item;
        }
    }
    panic("UID not found " + uid)
}

func OverrideForProvince(country string, province string) OverrideData {
    for _, item := range(overridedata) {
        if item.Country == country && item.Province == province && item.County == "" {
            return item;
        }
    }
    panic("item not found: " + country + ", " + province)
}

func OverrideForFips(fips string) OverrideData {
    if fips == "" { return OverrideData{} }
    for _, item := range(overridedata) {
        if item.Fips == fips {
            return item;
        }
    }
    panic("item not found: (fips)" + fips)
}
