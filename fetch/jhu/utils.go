package jhu

import (
	"covid-tracker/utils"
	"fmt"
	"io"
	"time"
    "strings"

	"github.com/pariz/gountries"
)

var CountryQuery = gountries.New()

func getDates(dateHeaders []string) []time.Time {
    dates := make([]time.Time, 0)
    for _, str := range dateHeaders {
        date, err := time.Parse(timeLayout, str)
        if err != nil {
            panic(err.Error())
        }
        dates = append(dates, date)
    }
    return dates
}

func getValuesForKind(kind string, value int) (confirmed int, deaths int) {
	if kind == "confirmed" {
		confirmed = value
	} else if kind == "deaths" {
		deaths = value
	} else {
        panic("Invalid kind passed to getValues: " + kind)
	}
	return confirmed, deaths
}

func fetchFromRepo(branch string, path string) io.ReadCloser {
    url := `https://raw.githubusercontent.com/CSSEGISandData/COVID-19/` + branch + "/" + path
    res, err := utils.Fetch(url)
    if err != nil { panic(err.Error()) }
    return res
}

// returns the index of the matched date or -1
func findMatchingValue(slice []TimeValue, value TimeValue) int {
	for i, current := range slice {
		if value.Date.Equal(current.Date) && value.Key == current.Key {
			return i
		}
	}
	return -1
}

func mergeSingleValue(values []TimeValue, value TimeValue) []TimeValue {
	index := findMatchingValue(values, value)
	if index != -1 {
		last := values[index]
		last.Confirmed += value.Confirmed
		last.Deaths += value.Deaths
		last.Recovered += value.Recovered
		values[index] = last
		return values
	} else {
		return append(values, value)
	}
}

type Key struct {
    Timestamp int64
    Key string
}

func consolidateValues(values []TimeValue) []TimeValue {
    mapper := make(map[Key]TimeValue)
    for _, value := range values {
        key := Key{
            Timestamp: value.Date.Unix(),
            Key: value.Key,
        }
        current, ok := mapper[key]
        if ok {
            current.Confirmed += value.Confirmed
            current.Deaths += value.Deaths
            current.Recovered += value.Recovered
            mapper[key] = current
        } else {
            mapper[key] = value
        }
    }

    result := make([]TimeValue, len(mapper))
    i := 0
	for _, value := range mapper {
        result[i] = value
        i++
	}

	return result
}

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
    "Sao Tome and Principe": "São Tomé and Príncipe",
}

func normalizeCountry(country string, province string) (string, string, string) {
    countryCode := ""
    if isCruise(country) {
        province = country
        country = "Cruise"
    } else {
        override, ok := countryOverrides[country]
        if ok {
            country = override
        }

        c, err := CountryQuery.FindCountryByName(country)
        if err != nil {
            fmt.Println(country)
            panic("ERROR: Country Not Found")
        }
        country = c.Name.Common
        countryCode = c.Codes.CCN3
    }
    return country, province, countryCode
}

func padFips(fips string) string {
    if len(fips) < 5 {
        return strings.Repeat("0", 5 - len(fips)) + fips
    }
    return fips
}
