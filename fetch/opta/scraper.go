package opta

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ignoreCounty(s string, c string) bool {
	return c == "Unassigned" ||
		c == "Unknown" ||
		s == "Diamond Princess" ||
		s == "Grand Princess" ||
		c == "Non-Utah resident" ||
		c == "Out-of-state" ||
		c == "Out of State"
}

var overrides = map[string]map[string]string{
	"AK": map[string]string{
		"Seward":      "Kenai Peninsula",
		"Homer":       "Kenai Peninsula",
		"Soldotna":    "Kenai Peninsula",
		"Sterling":    "Kenai Peninsula",
		"Palmer":      "Anchorage",
		"Eagle River": "Anchorage",
		"North Pole":  "Fairbanks North Star",
		"Gridwood":    "Anchorage",
	},
	"IL": map[string]string{
		"La Salle": "LaSalle",
	},
	"IN": map[string]string{
		"Verm.": "Vermillion",
	},
	"KY": map[string]string{
		"Northern Kentucky": "Unassigned",
	},
	"MA": map[string]string{
		"Dukes and Nantucket": "Dukes",
	},
	"MI": map[string]string{
		"Wayne--Detroit": "Wayne",
		"Other":          "Unassigned",
	},
	"MO": map[string]string{
		"Kansas City": "Jackson",
		"Joplin":      "Jasper",
		"TBD":         "Unassigned",
		"Phelps Maries":         "Phelps",
	},
	"MN": map[string]string{
		"Filmore": "Fillmore",
	},
	"OK": map[string]string{
		"Out-of-State": "Unassigned",
	},
	"PR": map[string]string{
		"Puerto Rico": "Unassigned",
	},
	"TN": map[string]string{
		"DeSoto": "Unassigned",
	},
	"TX": map[string]string{
		"De Witt": "DeWitt",
	},
	"UT": map[string]string{
		"Weber-Morgan":   "Weber",
		"Southwest Utah": "Unassigned",
		"TriCounty":      "Uintah",
		"Grant":          "Grand",
		"Unitah":         "Uintah",
		"Unassigned Southwest":         "Unassigned",
	},
	"VI": map[string]string{
		"St. Croix":  "St. Croix Island",
		"St. John":   "St. John Island",
		"St. Thomas": "St. Thomas Island",
	},
	"VA": map[string]string{
		"Virginia Beach":  "Virginia Beach City",
		"Alexandria":      "Alexandria City",
		"Harrisonburg":    "Harrisonburg City",
		"Charlottesville": "Charlottesville City",
		"Williamsburg":    "Williamsburg City",
		"Norfolk":         "Norfolk City",
		"Portsmouth":      "Portsmouth City",
		"Suffolk":         "Suffolk City",
		"Newport News":    "Newport News City",
		"Chesapeake":      "Chesapeake City",
		"Danville":        "Danville City",
		"Radford":         "Radford City",
		"Lynchburg":       "Lynchburg City",
		"Fredericksburg":  "Fredericksburg City",
		"Hampton":         "Hampton City",
		"Poquoson":        "Poquoson City",
		"Galax":           "Galax City",
		"Bristol":         "Bristol City",
		"Hopewell":        "Hopewell City",
		"Winchester":      "Winchester City",
		"Manassas Park":   "Manassas Park City",
	},
}

type optaItem struct {
	Id            int    `json:"id"`
	ConfirmedDate string `json:"confirmed_date"`
	PeopleCount   int    `json:"people_count"`
	DeathCount    int    `json:"die_count"`
	Comments      string `json:"comments_en"`
	State         string `json:"state_name"`
	County        string `json:"county"`
	Num           int    `json:"num"`
}

const MAIN_URL = "https://coronavirus.1point3acres.com"
const JS_PREFIX = "/_next/static/"
const TIME_LAYOUT = "1/2/2006"

var fipsData = ParseFips()

func mainPageHtml() (string, error) {
	return utils.FetchString(MAIN_URL)
}

func jsChunks(html string) []string {
	jsRegexp, _ := regexp.Compile(`chunks[^"]+\.js`)
	return jsRegexp.FindAllString(html, -1)
}

func chunkUrls(chunks []string) []string {
	for i, chunk := range chunks {
		chunks[i] = MAIN_URL + JS_PREFIX + chunk
	}
	return chunks
}

func fetchStrings(urls []string) ([]string, error) {
	data := make([]string, len(urls))

	for i, url := range urls {
		result, err := utils.FetchString(url)
		if err != nil {
			return nil, err
		}
		data[i] = result
	}
	return data, nil
}

func isValidJsData(data string) bool {
	return strings.Contains(data, "Snohomish")
}

func filterValidData(jsDatas []string) string {
	for _, jsData := range jsDatas {
		if isValidJsData(jsData) {
			return jsData
		}
	}

	fmt.Println("Found no valid data")
	return ""
}

func cleanData(jsData string) string {
	hexUnicode, err := regexp.Compile(`\\x(..)`)
	if err != nil {
		panic("Invalid regexp")
	}

	str := strings.ReplaceAll(jsData, "\\'", "'")
	str = strings.ReplaceAll(str, `\\"`, `\"`)
	str = hexUnicode.ReplaceAllString(str, `\u00$1`)
	return str
}

func toJsonString(jsData string) (string, string, error) {
	confirmedRes := strings.Split(jsData, "JSON.parse('")[3]
	confirmed := strings.Split(confirmedRes, "')}")[0]

	deathsRes := strings.Split(jsData, "JSON.parse('")[5]
	deaths := strings.Split(deathsRes, "')}")[0]

	return cleanData(confirmed), cleanData(deaths), nil
}

func toOpta(jsonString string) ([]optaItem, error) {
	var optas []optaItem

	unquoted, err := strconv.Unquote("`" + jsonString + "`")
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(strings.NewReader(unquoted))
	err = dec.Decode(&optas)
	if err != nil {
		return nil, err
	}

	return optas, nil
}

func checkList(optaList []optaItem) bool {
	for _, item := range optaList {
		result := item.County == "Snohomish"
		if result {
			return result
		}
	}
	return false
}

func convertItem(item optaItem, i int, length int) (data.CountyData, error) {
	date, err := time.Parse(TIME_LAYOUT, item.ConfirmedDate+"/2020")
	if err != nil {
		return data.CountyData{}, err
	}
	county := strings.ReplaceAll(item.County, "\u200b", "")
	stateOverride, ok := overrides[item.State]
	if ok {
		countyOverride, ok := stateOverride[county]
		if ok {
			county = countyOverride
		}
	}

	countyKey := item.State + "-" + county
	fips, ok := fipsData[strings.ToLower(countyKey)]

	if !ok && !ignoreCounty(item.State, county) {
		fmt.Println(strconv.Itoa(i) + ", " + strconv.Itoa(length))
		panic("County not found: " + strconv.Quote(countyKey))
	}

	return data.CountyData{
		ExternalId: strconv.Itoa(item.Id),
		StateCode:  item.State,
		State:      utils.StateCodes[item.State],
		County:     fips.Name,
		Confirmed:  item.PeopleCount,
		Deaths:     item.DeathCount,
		Date:       date,
		CountyKey:  countyKey,
		FipsId:     fips.Fips,
	}, nil
}

func convertItems(items []optaItem) ([]data.CountyData, error) {
	result := make([]data.CountyData, len(items))

	for i, item := range items {
		value, err := convertItem(item, i, len(items))
		if err != nil {
			return nil, err
		}

		result[i] = value
	}

	return result, nil
}

func hasTime(times []time.Time, t time.Time) bool {
	for _, t2 := range times {
		if t2.Equal(t) {
			return true
		}
	}
	return false
}

func hasString(strs []string, t string) bool {
	for _, t2 := range strs {
		if t2 == t {
			return true
		}
	}
	return false
}

func collectDates(data []data.CountyData) []time.Time {
	result := make([]time.Time, 0)

	for _, item := range data {
		if !hasTime(result, item.Date) {
			result = append(result, item.Date)
		}
	}
	return result
}

func collectCountyKeys(items []data.CountyData) map[string]data.CountyData {
	result := make(map[string]data.CountyData)

	for _, item := range items {
		if _, ok := result[item.CountyKey]; !ok {
			result[item.CountyKey] = data.CountyData{
				ExternalId: "",
				StateCode:  item.StateCode,
				State:      item.State,
				County:     item.County,
				Confirmed:  0,
				Deaths:     0,
				CountyKey:  item.CountyKey,
				FipsId:     item.FipsId,
			}
		}
	}
	return result
}

type timeSlice []time.Time

func (a timeSlice) Len() int           { return len(a) }
func (a timeSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a timeSlice) Less(i, j int) bool { return a[i].Before(a[j]) }

func mapData(items []data.CountyData) map[time.Time]map[string]data.CountyData {
	result := make(map[time.Time]map[string]data.CountyData)
	for _, item := range items {
		time, ok := result[item.Date]
		if !ok {
			time = make(map[string]data.CountyData)
		}
		key, ok := time[item.CountyKey]
		if !ok {
			key = item
		} else {
			key.Confirmed += item.Confirmed
			key.Deaths += item.Deaths
			if key.ExternalId == "" {
				key.ExternalId = item.ExternalId
			}
		}

		time[item.CountyKey] = key
		result[item.Date] = time
	}
	return result
}

func GetData() ([]data.CountyData, error) {
	html, err := mainPageHtml()
	if err != nil {
		return nil, err
	}

	urls := chunkUrls(jsChunks(html))
	jsData, err := fetchStrings(urls)
	if err != nil {
		return nil, err
	}

	confirmedString, deathsString, err := toJsonString(filterValidData(jsData))
	if err != nil {
		return nil, err
	}

	confirmedData, err := toOpta(confirmedString)
	if err != nil {
		return nil, err
	}

	deathsData, err := toOpta(deathsString)
	if err != nil {
		return nil, err
	}

	optaData := append(confirmedData, deathsData...)

	result, err := convertItems(optaData)
	if err != nil {
		return nil, err
	}

	m := mapData(result)

	allDates := timeSlice(collectDates(result))
	min := allDates[0]
	max := allDates[len(allDates)-1]
	sort.Sort(allDates)
	keys := collectCountyKeys(result)

	runningConfirmed := make(map[string]int)
	runningDeaths := make(map[string]int)

	for key, _ := range keys {
		runningConfirmed[key] = 0
		runningDeaths[key] = 0
	}

	newresult := make([]data.CountyData, 0)
	date := min
	for !date.After(max) {
		keyMap := m[date]
		for key, base := range keys {
			val, ok := keyMap[key]
			if ok {
				runningConfirmed[key] += val.Confirmed
				runningDeaths[key] += val.Deaths
				val.Confirmed = runningConfirmed[key]
				val.Deaths = runningDeaths[key]
				newresult = append(newresult, val)
			} else {
				newresult = append(newresult, data.CountyData{
					ExternalId: "",
					StateCode:  base.StateCode,
					State:      base.State,
					County:     base.County,
					Confirmed:  runningConfirmed[key],
					Deaths:     runningDeaths[key],
					Date:       date,
					CountyKey:  base.CountyKey,
					FipsId:     base.FipsId,
				})
			}
		}

		date = date.AddDate(0, 0, 1)
	}

	return newresult, nil
}
