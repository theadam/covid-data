package jhu

import (
	"covid-tracker/utils"
	"covid-tracker/data"
	"io"
    "time"
	"encoding/csv"
	"strings"
	"strconv"
	"errors"
    "fmt"
    "github.com/pariz/gountries"
)

const timeLayout = "1/2/06"

func isCruise(country string) bool {
    return country == "Diamond Princess" || country == "MS Zaandam"
}

var countryOverrides = map[string]string{
    "Cabo Verde": "Cape Verde",
    "Congo (Brazzaville)": "Democratic Republic of the Congo",
    "Congo (Kinshasa)": "Democratic Republic of the Congo",
    "Cote d'Ivoire": "Ivory Coast",
    "Czechia": "Czech Republic",
    "Eswatini": "Swaziland",
    "Holy See": "Vatican City",
    "Korea, South": "South Korea",
    "North Macedonia": "Macedonia",
    "Taiwan*": "Taiwan",
    "US": "United States",
    "The West Bank and Gaza": "Palestine",
    "West Bank and Gaza": "Palestine",
    "Kosovo": "Republic of Serbia",
    "Burma": "Republic of the Union of Myanmar",
}


type rawValue struct {
	state     string
	country   string
	lat       string
	long      string
	main      bool
	date      time.Time
	confirmed int
	deaths    int
}

type DateValue struct {
	Date      time.Time
	Confirmed int
	Deaths    int
}

type StateData struct {
	State   string
	Country string
	Data    []DateValue
	Lat     string
	Long    string
}

type StateMap map[string]StateData

type CountryData struct {
	States  StateMap
	Data    []DateValue
	Lat     string
	Long    string
	Country string
}

type covidData map[string]CountryData

func getValues(kind string, value int) (confirmed int, deaths int, err error) {
	if kind == "confirmed" {
		confirmed = value
	} else if kind == "deaths" {
		deaths = value
	} else {
		err = errors.New(fmt.Sprintf("Invalid kind: %s", kind))
	}
	return confirmed, deaths, err
}


func fetchFromRepo(path string) (io.ReadCloser, error) {
	url := `https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/` + path
	return utils.Fetch(url)
}

func fetchTimeSeries(path string, kind string) ([]data.DataPoint, error) {
	body, err := fetchFromRepo(path)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	reader := csv.NewReader(body)

	points := make([]data.DataPoint, 0)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	dateHeaders := header[4:]

    dates, err := getDates(dateHeaders)
	if err != nil {
		return nil, err
	}

	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return points, err
		}

        state, country, code := getFields(records)

		dateData := records[4:]

		// parse each item in the time series
		for i, item := range dateData {
            if item == "" {
                item = "0"
            }
			number, err := strconv.Atoi(item)
			if err != nil {
				return points, err
			}

			confirmed, deaths, err := getValues(kind, number)

			if err != nil {
				return nil, err
			}

            date := dates[i]

			points = append(points, data.DataPoint{
				Confirmed: confirmed,
				Deaths:    deaths,
				Date:      date,
				Country:   country,
                CountryCode: code,
                State:     state,
                ExternalState: records[0],
                ExternalCountry: records[1],
				Lat:       records[2],
				Long:      records[3],
			})
		}
	}

	return points, nil
}

func fetchConfirmed() ([]data.DataPoint, error) {
    return fetchTimeSeries("csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv", "confirmed")
}

func fetchDeaths() ([]data.DataPoint, error) {
    return fetchTimeSeries("csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv", "deaths")
}

func getFields(records []string) (state string, country string, code string) {
	state = records[0]
	country = records[1]
    if isCruise(country) {
        return country, "Cruise", ""
    }

    override, ok := countryOverrides[country]
    if ok {
        country = override
    }

    query := gountries.New()
    c, err := query.FindCountryByName(country)
    if err != nil {
        fmt.Println(country)
        panic("ERROR: Country Not Found")
    }
    country = c.Name.Common
    code = c.Codes.CCN3
    if country == "United States" && strings.Contains(state, ", ") {

		items := strings.Split(state, ", ")
		result, ok := utils.StateCodes[strings.TrimSpace(items[1])]

		if ok {
			state = result
		}
	}
    return state, country, code
}

func getDates(dateHeaders []string) (dates []time.Time, err error) {
    dates = make([]time.Time, 0)
    for _, str := range dateHeaders {
        date, err := time.Parse(timeLayout, str)
        if err != nil {
            return nil, err
        }
        dates = append(dates, date)
    }
    return dates, nil
}


// returns the index of the matched date or -1
func findMatching(slice []data.DataPoint, point data.DataPoint) int {
	for i, value := range slice {
		if value.Date.Equal(point.Date) && value.State == point.State && value.Country == point.Country {
			return i
		}
	}
	return -1
}

func updateData(points []data.DataPoint, point data.DataPoint) []data.DataPoint {
	index := findMatching(points, point)
	if index != -1 {
		last := points[index]
		last.Confirmed += point.Confirmed
		last.Deaths += point.Deaths
		points[index] = last
		return points
	} else {
		return append(points, point)
	}
}

func fetchAll() (masterList []data.DataPoint, err error) {
    confirmed, err := fetchConfirmed()
	if err != nil { return nil, err }

    deaths, err := fetchDeaths()
	if err != nil { return nil, err }

	master := append(confirmed, deaths...)

    return master, nil
}

func consolidatePoints(points []data.DataPoint) []data.DataPoint {
	result := make([]data.DataPoint, 0)

	for _, point := range points {
        result = updateData(result, point);
	}

	return result
}

func GetData(start time.Time) ([]data.DataPoint, error) {
    master, err := fetchAll()
    if err != nil { return nil, err }

    master = utils.FilterDataPointByDate(master, start)

    return consolidatePoints(master), nil
}
