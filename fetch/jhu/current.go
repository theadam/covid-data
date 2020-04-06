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

func fetchCurrent(globalDate time.Time, usDate time.Time) ([]data.DataPoint, []data.CountyData) {
	body := fetchFromRepo("web-data", "data/cases.csv")
	defer body.Close()

	reader := csv.NewReader(body)

	countyData := make([]data.CountyData, 0)
	worldData := make([]data.DataPoint, 0)
    // ignore header
    _, err := reader.Read()
	if err != nil { panic(err.Error()) }

	for {
		columns, err := reader.Read()
		if err == io.EOF { break }

        county := columns[1]
        country := columns[3]
        lat := columns[5]
        long := columns[6]

        confirmed, err := strconv.Atoi(columns[7])
        if err != nil { panic(err.Error()) }

        deaths, err := strconv.Atoi(columns[8])
        if err != nil { panic(err.Error()) }

        if county != "" {
            if country != "US" {
                fmt.Println(strings.Join(columns, ", "))
                panic("County found not in the US")
            }
            fipsData := FipsMap[columns[0]]
            countyData = append(countyData, data.CountyData{
                FipsId: fipsData.Fips,
                State: utils.StateCodes[fipsData.StateCode],
                StateCode: fipsData.StateCode,
                County: fipsData.Name,
                Confirmed: confirmed,
                Deaths: deaths,
                Date: usDate,
                Lat: lat,
                Long: long,
            })

        } else {
            country, province, countryCode := normalizeCountry(country, columns[2])

            worldData = append(worldData, data.DataPoint{
                Province: province,
                Country: country,
                CountryCode: countryCode,
                Confirmed: confirmed,
                Deaths: deaths,
                Date: globalDate,
                Lat: lat,
                Long: long,
                ExternalCountry: columns[1],
            })
        }
    }

    return worldData, countyData
}
