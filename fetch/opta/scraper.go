package opta

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type optaItem struct {
    Id int `json:"id"`
    ConfirmedDate string `json:"confirmed_date"`
    PeopleCount int `json:"people_count"`
    DeathCount int `json:"die_count"`
    Comments string `json:"comments_en"`
    State string `json:"state_name"`
    County string `json:"county"`
    Num int `json:"num"`
}

const MAIN_URL = "https://coronavirus.1point3acres.com"
const JS_PREFIX = "/_next/static/"
const TIME_LAYOUT = "1/2/2006"

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

func toJsonString(jsData string) (string, error) {
    hexUnicode, err := regexp.Compile(`\\x(..)`)
    if err != nil { return "", err }
    res := strings.Split(jsData, "JSON.parse('")[3]
    str := strings.Split(res, "')}")[0]
    str = strings.ReplaceAll(str, "\\'", "'")
    str = strings.ReplaceAll(str, `\\"`, `\"`)
    str = hexUnicode.ReplaceAllString(str, `\u00$1`)
    return str, nil
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

func convertItem(item optaItem) (data.CountyData, error) {
    date, err := time.Parse(TIME_LAYOUT, item.ConfirmedDate + "/2020")
    if err != nil {
        return data.CountyData{}, err
    }
    return data.CountyData{
        ExternalId: strconv.Itoa(item.Id),
        StateCode: item.State,
        State: utils.StateCodes[item.State],
        County: item.County,
        Confirmed: item.PeopleCount,
        Deaths: item.DeathCount,
        Date: date,
    }, nil
}

func convertItems(items []optaItem) ([]data.CountyData, error) {
    result := make([]data.CountyData, len(items))

    for i, item := range items {
        value, err := convertItem(item)
        if err != nil {
            return nil, err
        }

        result[i] = value
    }

    return result, nil
}

func GetData() ([]data.CountyData, error) {
    html, err :=  mainPageHtml()
    if err != nil { return nil, err }

    urls := chunkUrls(jsChunks(html))
    jsData, err := fetchStrings(urls)
    if err != nil { return nil, err }

    jsonString, err := toJsonString(filterValidData(jsData))
    if err != nil { return nil, err }

    optaData, err := toOpta(jsonString)
    if err != nil { return nil, err }

    result, err := convertItems(optaData)
    if err != nil { return nil, err }

    return result, nil
}
