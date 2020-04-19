package jhu

import (
	"covid-tracker/data"
	"encoding/csv"
	"io"
	"strconv"
)

const timeLayout = "1/2/06"

type DataRow interface {
	TimeColumns([]string, string) []string
	ToDataPoint([]string) *data.DataPoint
}

func parseTimeSeries(
	path string,
	kind string,
	dataRow DataRow,
) []*data.DataPoint {
	body := fetchFromRepo("master", path)
	defer body.Close()

	reader := csv.NewReader(body)

	rows := make([]*data.DataPoint, 0)
	header, err := reader.Read()
	if err != nil {
		panic(err.Error())
	}

	dateHeaders := dataRow.TimeColumns(header, kind)
	dates := getDates(dateHeaders)

	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		dateData := dataRow.TimeColumns(columns, kind)

		// parse each item in the time series
		for i, item := range dateData {
			date := dates[i]

			if item == "" {
				item = "0"
			}
            float, err := strconv.ParseFloat(item, 64)
			if err != nil {
				panic(err.Error())
			}
            number := int(float)

			dataPoint := dataRow.ToDataPoint(columns)

			if dataPoint != nil {
				confirmed, deaths := getValuesForKind(kind, number)
                dataPoint.Date = date
                dataPoint.Confirmed = confirmed
                dataPoint.Deaths = deaths

				rows = append(rows, dataPoint)
			}
		}
	}
	return rows
}

func fetchGlobalConfirmed() []*data.DataPoint {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv",
		"confirmed",
		&GlobalTS,
	)
}

func fetchGlobalDeaths() []*data.DataPoint {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv",
		"deaths",
		&GlobalTS,
	)
}

func fetchUsConfirmed() []*data.DataPoint{
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv",
		"confirmed",
		&UsTS,
	)
}

func fetchUsDeaths() []*data.DataPoint {
	return parseTimeSeries(
		"csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_US.csv",
		"deaths",
		&UsTS,
	)
}

func timeSeriesGlobals() []*data.DataPoint {
	confirmed := fetchGlobalConfirmed()
	deaths := fetchGlobalDeaths()
	return consolidateValues(append(confirmed, deaths...))
}

func timeSeriesUs() []*data.DataPoint {
	confirmed := fetchUsConfirmed()
	deaths := fetchUsDeaths()
	return consolidateValues(append(confirmed, deaths...))
}
