package jhu

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

const timeLayout = "1/2/06"

type TimeValue struct {
	Date      time.Time
	Confirmed int
	Deaths    int
	Recovered int
	ExtraData map[string]string
	Key       string
}

type TimeseriesRow interface {
	TimeColumns([]string, string) []string
	ExtractExtraData([]string) map[string]string
	Key(map[string]string) string
	Skip(map[string]string) bool
}

func parseTimeSeries(
	path string,
	kind string,
	start time.Time,
	timeSeriesRow TimeseriesRow,
) []TimeValue {
	body := fetchFromRepo("master", path)
	defer body.Close()

	reader := csv.NewReader(body)

	rows := make([]TimeValue, 0)
	header, err := reader.Read()
	if err != nil {
		panic(err.Error())
	}

	dateHeaders := timeSeriesRow.TimeColumns(header, kind)
	dates := getDates(dateHeaders)

	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		dateData := timeSeriesRow.TimeColumns(columns, kind)

		// parse each item in the time series
		for i, item := range dateData {
			date := dates[i]
			// Only handle dates after start
			if date.Before(start) {
				continue
			}

			if item == "" {
				item = "0"
			}
            float, err := strconv.ParseFloat(item, 64)
			if err != nil {
				panic(err.Error())
			}
            number := int(float)

			extraData := timeSeriesRow.ExtractExtraData(columns)

			if !timeSeriesRow.Skip(extraData) {
				confirmed, deaths := getValuesForKind(kind, number)
				timeData := TimeValue{
					Date:      date,
					Confirmed: confirmed,
					Deaths:    deaths,
					ExtraData: extraData,
					Key:       timeSeriesRow.Key(extraData),
				}

				rows = append(rows, timeData)
			}
		}
	}
	return rows
}

func fetchGlobalConfirmed(start time.Time) []TimeValue {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv",
		"confirmed",
		start,
		&GlobalTS,
	)
}

func fetchGlobalDeaths(start time.Time) []TimeValue {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv",
		"deaths",
		start,
		&GlobalTS,
	)
}

func fetchUsTerritoriesConfirmed(start time.Time) []TimeValue {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv",
		"confirmed",
		start,
		&UsTerritoryTS,
	)
}

func fetchUsTerritoriesDeaths(start time.Time) []TimeValue {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_US.csv",
		"deaths",
		start,
		&UsTerritoryTS,
	)
}

func fetchUsConfirmed(start time.Time) []TimeValue {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv",
		"confirmed",
		start,
		&UsTS,
	)
}

func fetchUsDeaths(start time.Time) []TimeValue {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_US.csv",
		"deaths",
		start,
		&UsTS,
	)
}

func timeSeriesGlobals(start time.Time) []TimeValue {
	confirmed := fetchGlobalConfirmed(start)
	deaths := fetchGlobalDeaths(start)
	return consolidateValues(append(confirmed, deaths...))
}

func timeSeriesUsTerritories(start time.Time) []TimeValue {
	confirmed := fetchUsTerritoriesConfirmed(start)
	deaths := fetchUsTerritoriesDeaths(start)
	return consolidateValues(append(confirmed, deaths...))
}

func timeSeriesUs(start time.Time) []TimeValue {
	confirmed := fetchUsConfirmed(start)
	deaths := fetchUsDeaths(start)
	return consolidateValues(append(confirmed, deaths...))
}
