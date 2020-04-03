package opta

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ignoreCounty(s string, c string) bool {
	return c == "Unassigned" ||
		c == "unassigned" ||
		c == "Unknown" ||
		s == "Diamond Princess" ||
		s == "Grand Princess" ||
		c == "Non-Utah resident" ||
		c == "Out-of-state" ||
		c == "Out of State" ||
		c == "Out of state" ||
		s == "NN" || (s == "GA" && c == "Chambers")
}

type stateCounty struct {
	State  string
	County string
}

var coords = map[[2]byte]string{{1, 1}: "one one", {2, 1}: "two one"}

var fullOverride = map[[2]string][2]string{
	{"NN", "Navajo, AZ"}: {"AZ", "Navajo"},
}

var overrides = map[string]map[string]string{
	"AK": map[string]string{
        //		"Kenai":       "Kenai Peninsula",
        //		"Seward":      "Kenai Peninsula",
        //		"Homer":       "Kenai Peninsula",
        //		"Soldotna":    "Kenai Peninsula",
        //		"Sterling":    "Kenai Peninsula",
        //		"Palmer":      "Anchorage",
        //		"Eagle River": "Anchorage",
        //		"Anchorage--Eagle River": "Anchorage",
        //		"North Pole":  "Fairbanks North Star",
        //		"Gridwood":    "Anchorage",
        //		"Wasilla":     "Matanuska-Susitna",
	},
	"IA": map[string]string{
        "Obrien": "O'Brien",
	},
	"ID": map[string]string{
        "Adam": "Adams",
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
	"LA": map[string]string{
		"Parish Under Investigation": "Unassigned",
	},
	"MA": map[string]string{
		"Dukes and Nantucket": "Dukes",
	},
	"MI": map[string]string{
		"Other":          "Unassigned",
	},
	"MO": map[string]string{
		"Kansas City":   "Jackson",
		"Joplin":        "Jasper",
		"TBD":           "Unassigned",
		"Phelps Maries": "Phelps",
	},
	"MN": map[string]string{
		"Filmore": "Fillmore",
	},
	"NH": map[string]string{
		"Hillsborough-other":     "Hillsborough",
		"Hillsborough-Manchester":     "Hillsborough",
		"Hillsborough-Nashua":     "Hillsborough",
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
		"De Witt":             "DeWitt",
		"Harris--Non Houston": "Harris",
		"Harris--Houston":     "Harris",
		"El Paso--Fort Bliss": "El Paso",
	},
	"UT": map[string]string{
		"Weber-Morgan":         "Weber",
		"Southwest Utah":       "Unassigned",
		"Central Utah":         "Unassigned",
		"TriCounty":            "Uintah",
		"Grant":                "Grand",
		"Unitah":               "Uintah",
		"Unassigned Southwest": "Unassigned",
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
		"Petersburg":      "Petersburg City",
		"Waynesboro":      "Waynesboro City",
		"Charke":          "Clarke",
		"Covington":       "Covington City",
		"Emporia":       "Emporia City",
		"Lexington":       "Lexington City",
		"Salem":          "Roanoke",
	},
}

type preOptaItem struct {
	Id             int    `json:"id"`
	ConfirmedDate  string `json:"confirmed_date"`
	PeopleCount    int    `json:"people_count"`
	DeathCount     int    `json:"die_count"`
	Comments       string `json:"comments_en"`
	State          string `json:"state_name"`
	County         string `json:"county"`
	RecoveredCount int    `json:"cured_count"`
	Num            int    `json:"num"`
}

type optaItem struct {
	Id             int    `json:"id"`
	ConfirmedDate  string `json:"confirmed_date"`
	PeopleCount    int    `json:"people_count"`
	DeathCount     int    `json:"die_count"`
	Comments       string `json:"comments_en"`
	State          string `json:"state_name"`
	County         string `json:"county"`
	RecoveredCount int    `json:"cured_count"`
	Num            int    `json:"num"`
	Orig           string
}

func (o *optaItem) UnmarshalJSON(data []byte) error {
	var v preOptaItem
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	o.Id = v.Id
	o.ConfirmedDate = v.ConfirmedDate
	o.PeopleCount = 0
	o.DeathCount = v.DeathCount
	o.RecoveredCount = v.DeathCount
	o.Comments = v.Comments
	o.State = v.State
	o.County = v.County
	o.Num = v.Num
	o.Orig = string(data)
	return nil
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

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func chunkUrls(chunks []string) []string {
	for i, chunk := range chunks {
		chunks[i] = MAIN_URL + JS_PREFIX + chunk
	}
	return unique(chunks)
}

func fetchStrings(urls []string) []string {
	data := make([]string, len(urls))

	for i, url := range urls {
		result, err := utils.FetchString(url)
		if err != nil {
			panic(err.Error())
		}
		data[i] = result
	}
	return data
}

func Unmarshal(jsonString string, v interface{}) error {
	data := cleanData(jsonString)
	cleaned := strings.ReplaceAll(data, `"`, `\u0022`)
	unquoted, err := strconv.Unquote(`"` + cleaned + `"`)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(unquoted), v)
}

func prettyJson(jsonString string) string {
	var v interface{}
	err := Unmarshal(jsonString, &v)
	if err != nil {
		fmt.Println(jsonString)
		panic("Failed to unmarshal pretty json: " + err.Error())
	}
	res, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(jsonString)
		panic("Failed to marshal pretty json: " + err.Error())
	}
	return string(res)
}

func auditAllJson(content string) {
	pieces := strings.Split(content, "JSON.parse('")[1:]
	for i, piece := range pieces {
		jsonString := strings.Split(piece, "')")[0]
		var result interface{}
		err := Unmarshal(jsonString, &result)
		if err != nil {
			panic(err.Error())
		}

		var item interface{}
		switch v := reflect.ValueOf(result); v.Kind() {
		case reflect.Map:
			{
				m := make(map[string]interface{})
				iter := v.MapRange()
				iter.Next()
				m[iter.Key().String()] = iter.Value().Interface()
				item = m
			}
		case reflect.Slice:
			{
				item = v.Index(0).Interface()
			}
		default:
			{
				fmt.Println(prettyJson(jsonString))
				panic("Unahndled kind: " + v.Kind().String())
			}
		}
		res, err := json.Marshal(item)
		if err != nil {
			fmt.Println(prettyJson(jsonString))
			panic(err.Error())
		}
		fmt.Println("\t" + strconv.Itoa(i) + ": " + string(res))
		fmt.Println()
	}
}

func auditJs(urls []string, contents []string) {
	for i, url := range urls {
		content := contents[i]
		if isValidJsData(content) {
			fmt.Println("URL: " + url)
			auditAllJson(content)
			fmt.Println()
		}
	}
	panic("Auditing")
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

func toJsonString(jsData string) (string, []string) {
	pieces := strings.Split(jsData, "JSON.parse('")[1:]

	confirmedRes := pieces[6]
	confirmed := strings.Split(confirmedRes, "')}")[0]

	deathsRes := pieces[3]
	deaths := strings.Split(deathsRes, "')}")[0]

	curedRes := pieces[7]
	recovered := strings.Split(curedRes, "')}")[0]

	return confirmed, []string{deaths, recovered}
}

func getFirstString(slice interface{}) string {
	first := reflect.ValueOf(slice).Index(0).Interface().(string)
	return first
}

func convertConfirmed(jsonString string) []optaItem {
	dateRegexp, _ := regexp.Compile(`\d{1,2}/\d{1,2}`)
	var data []map[string]interface{}
	err := Unmarshal(jsonString, &data)
	if err != nil {
		panic("Failed to parse confirmed: " + err.Error())
	}

	results := make([]optaItem, 0)
	for _, item := range data {
		orig, _ := json.Marshal(item)
		state := getFirstString(item["state_name"])
		county := getFirstString(item["county"])

		for key, val := range item {
			if dateRegexp.MatchString(key) {
				results = append(results, optaItem{
					State:         state,
					County:        county,
					PeopleCount:   int(reflect.ValueOf(val).Float()),
					ConfirmedDate: key,
					Orig:          string(orig),
				})
			}
		}
	}
	return results
}

func toOpta(jsonString string) []optaItem {
	var optas []optaItem
	err := Unmarshal(jsonString, &optas)
	if err != nil {
		panic(err.Error())
	}
	return optas
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
		fmt.Println(item.Orig)
		return data.CountyData{}, err
	}
	county := strings.ReplaceAll(item.County, "\u200b", "")
	state := item.State

	override, ok := fullOverride[[2]string{state, county}]

	if ok {
		county = override[1]
		state = override[0]
	}

    if strings.Contains(county, "--") {
        county = strings.Split(county, "--")[0]
    }

    county = strings.TrimSpace(county)

	stateOverride, ok := overrides[state]
	if ok {
		countyOverride, ok := stateOverride[county]
		if ok {
			county = countyOverride
		}
	}

	countyKey := state + "-" + county
	fips, ok := fipsData[strings.ToLower(countyKey)]

	if !ok && !ignoreCounty(state, county) {
		fmt.Println(strconv.Itoa(i) + ", " + strconv.Itoa(length))
		fmt.Println(item.Orig)
		fmt.Println(item.ConfirmedDate)
		fmt.Println(item.Comments)
		panic("County not found: " + strconv.Quote(countyKey))
	}

	return data.CountyData{
		ExternalId: strconv.Itoa(item.Id),
		StateCode:  state,
		State:      utils.StateCodes[state],
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

func GetData(start time.Time) ([]data.CountyData, error) {
	html, err := mainPageHtml()
	if err != nil {
		return nil, err
	}

	urls := chunkUrls(jsChunks(html))
	jsData := fetchStrings(urls)

	confirmedString, rest := toJsonString(filterValidData(jsData))

	optaData := convertConfirmed(confirmedString)

	for _, d := range rest {
		optaData = append(optaData, toOpta(d)...)
	}

	result, err := convertItems(optaData)
	if err != nil {
		return nil, err
	}

	m := mapData(result)

	allDates := timeSlice(collectDates(result))
	sort.Sort(allDates)
	min := allDates[0]
	max := allDates[len(allDates)-1]
	keys := collectCountyKeys(result)

	runningConfirmed := make(map[string]int)
	runningDeaths := make(map[string]int)
	runningRecovered := make(map[string]int)

	for key, _ := range keys {
		runningConfirmed[key] = 0
		runningDeaths[key] = 0
		runningRecovered[key] = 0
	}

	newresult := make([]data.CountyData, 0)
	date := min

	for !date.After(max) {
		keyMap := m[date]
		for key, base := range keys {
			val, ok := keyMap[key]
			if ok {
                if val.Confirmed != 0 {
                    runningConfirmed[key] = val.Confirmed
                }
				runningDeaths[key] += val.Deaths
				runningRecovered[key] += val.Recovered
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
					Recovered:  runningRecovered[key],
					Date:       date,
					CountyKey:  base.CountyKey,
					FipsId:     base.FipsId,
				})
			}
		}

		date = date.AddDate(0, 0, 1)
	}

    newresult = utils.FilterCountyDataByDate(newresult, start)

	return newresult, nil
}
