package jhu

import (
	"covid-tracker/utils"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	. "github.com/ahmetb/go-linq"
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

func fetchOverrides() (map[string]OverrideData, map[string]OverrideData, map[string]OverrideData) {
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
            if country != "Cruise" && columns[5] != "Unassigned" && !strings.HasPrefix(columns[5], "Out of") && province != "Recovered" && !utils.IsOrganization(province){
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
    uidMap := make(map[string]OverrideData)
    fipsMap := make(map[string]OverrideData)
    provinceCountryMap := make(map[string]OverrideData)

    From(data).Where(func(inter interface{}) bool {
        item := inter.(OverrideData)
        return item.UID != ""
    }).ToMapBy(&uidMap, utils.Field("UID"), utils.Id)

    From(data).Where(func(inter interface{}) bool {
        item := inter.(OverrideData)
        return item.Fips != ""
    }).ToMapBy(&fipsMap, utils.Field("Fips"), utils.Id)

    From(data).Where(func(inter interface{}) bool {
        item := inter.(OverrideData)
        return item.County == ""
    }).ToMapBy(&provinceCountryMap, func(inter interface{}) interface{} {
        item := inter.(OverrideData)
        return item.Province + "-" + item.Country
    }, utils.Id)

    return uidMap, fipsMap, provinceCountryMap
}
var uidMap, fipsMap, provinceCountryMap = fetchOverrides()

func OverrideFromUID(uid string) OverrideData {
    val, ok := uidMap[uid]
    if !ok {
        panic("UID not found " + uid)
    }
    return val
}

func OverrideForFips(fips string) OverrideData {
    val, ok := fipsMap[fips]
    if !ok {
        panic("Fips not found: " + fips)
    }
    return val
}

func OverrideForProvince(country string, province string) OverrideData {
    val, ok := provinceCountryMap[province + "-" + country]
    if !ok {
        panic("Item not found: " + country + ", " + province)
    }
    return val
}
