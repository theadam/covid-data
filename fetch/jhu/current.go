package jhu

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func fetchCurrent(
	date time.Time,
	latestGlobals map[string]TimeValue,
	latestUs map[string]TimeValue,
) ([]data.DataPoint, []data.CountyData) {
	usedGlobals := make(map[string]bool)
	for k, _ := range latestGlobals {
		usedGlobals[k] = false
	}
	usedUs := make(map[string]bool)
	for k, _ := range latestUs {
		usedUs[k] = false
	}

	body := fetchFromRepo("web-data", "data/cases.csv")
	defer body.Close()

	reader := csv.NewReader(body)

	countyData := make([]data.CountyData, 0)
	worldData := make([]data.DataPoint, 0)
	// ignore header
	_, err := reader.Read()
	if err != nil {
		panic(err.Error())
	}

	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		}

		county := columns[9]
		country := columns[1]
		lat := columns[3]
		long := columns[4]

		confirmed, err := strconv.Atoi(columns[5])
		if err != nil {
			panic(err.Error())
		}

		deaths, err := strconv.Atoi(columns[6])
		if err != nil {
			panic(err.Error())
		}

		if county != "" {
			if country != "US" {
                fmt.Println(country)
				fmt.Println(strings.Join(columns, ", "))
				panic("County found not in the US")
			}
            state := columns[0]
			fipsData := FipsMap[padFips(columns[10])]
            if fipsData.Fips == "" {
                f := extractFips(columns[15])
                fipsData = FipsMap[f]
            }
            if fipsData.Fips != "" {
                state = utils.StateCodes[fipsData.StateCode]
                county = fipsData.Name
            }
            key := makeUsKey(fipsData.Fips, county, state)
			usedUs[key] = true
			countyData = append(countyData, data.CountyData{
				FipsId:    fipsData.Fips,
				State:     state,
				County:    county,
				Confirmed: confirmed,
				Deaths:    deaths,
				Date:      date,
				Lat:       lat,
				Long:      long,
                Population: OverrideForFips(fipsData.Fips).Population,
			})
		} else {
			country, province, countryCode := normalizeCountry(country, columns[0])
			if !skipGlobal(country, province) {
				key := makeGlobalKey(country, province)
				usedGlobals[key] = true
				worldData = append(worldData, data.DataPoint{
					Province:        province,
					Country:         country,
					CountryCode:     countryCode,
					Confirmed:       confirmed,
					Deaths:          deaths,
					Date:            date,
					Lat:             lat,
					Long:            long,
					ExternalCountry: columns[1],
                    Population: OverrideForProvince(country, province).Population,
				})
			}
		}
	}

	for k, v := range usedGlobals {
		if v == false {
			addedValue := toDataPoint(latestGlobals[k])
			addedValue.Date = date
			worldData = append(worldData, addedValue)
		}
	}

	for k, v := range usedUs {
		if v == false {
			addedValue := toCountyData(latestUs[k])
			addedValue.Date = date
			countyData = append(countyData, addedValue)
		}
	}
	return worldData, countyData
}
