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
	date *time.Time,
	latest map[data.LocationKey]*data.DataPoint,
) []*data.DataPoint {
	used := make(map[data.LocationKey]bool)
	for k, _ := range latest {
		used[k] = false
	}

	body := fetchFromRepo("web-data", "data/cases.csv")
	defer body.Close()

	reader := csv.NewReader(body)

	points := make([]*data.DataPoint, 0)
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
        country, province, countryCode := normalizeCountry(country, columns[0])
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

        fipsData := FipsMap[padFips(columns[10])]

        if fipsData.Fips == "" {
            f := extractFips(columns[15])
            fipsData = FipsMap[f]
        }

        var population int
        if fipsData.Fips != "" {
            province = utils.StateCodes[fipsData.StateCode]
            county = fipsData.Name
            population = OverrideForFips(fipsData.Fips).Population
        } else {
            population = OverrideForProvince(country, province).Population
        }

        if county != "" && country != "United States" {
            fmt.Println(country)
            fmt.Println(strings.Join(columns, ", "))
            panic("County found not in the US")
        }

        if !skipGlobal(country, province) {
            point := data.DataPoint{
                Country: country,
                CountryCode: countryCode,

                Province: province,

                County: county,
                FipsId: fipsData.Fips,

                Lat: lat,
                Long: long,

                Population: population,

                Date: *date,
                Confirmed: confirmed,
                Deaths: deaths,
            }
            used[*point.LocationKey()] = true
            points = append(points, &point)
        }
	}

	for k, v := range used {
		if v == false {
            l := latest[k]
            addedValue := data.DataPoint{
                Country: l.Country,
                CountryCode: l.CountryCode,

                Province: l.Province,

                County: l.County,
                FipsId: l.FipsId,

                Lat: l.Lat,
                Long: l.Long,

                Date: *date,
                Confirmed: l.Confirmed,
                Deaths: l.Deaths,
                Population: l.Population,
            }
			points = append(points, &addedValue)
		}
	}
	return points
}
