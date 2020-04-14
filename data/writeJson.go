package data

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

var base = "./client/src/data/"

type Date struct {
    time.Time
}

type DateData struct {
	Date      Date `json:"date"`
	Confirmed int       `json:"confirmed"`
	Deaths    int       `json:"deaths"`
}

const timeLayout = "2006-01-02"

func (u *Date) MarshalJSON() ([]byte, error) {
    return []byte(`"` + u.Format(timeLayout) + `"`), nil
}

func makeDateData(date time.Time, confirmed int, deaths int) DateData {
    return DateData{
        Date: Date{date},
        Confirmed: confirmed,
        Deaths: deaths,
    }
}

func writeFile(json string, file string) {
	path := base + file
	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()
	_, err = f.WriteString(json)
	if err != nil {
		panic(err.Error())
	}
}

func getDates(db *gorm.DB, model interface{}) (time.Time, time.Time) {
	var max time.Time
	var min time.Time
	db.Select("max(date) as max, min(date) as min").Model(model).Row().Scan(&max, &min)
	return min, max
}

func dateRange(db *gorm.DB) (time.Time, time.Time) {
	min, max := getDates(db, &Point)
	cmin, cmax := getDates(db, &CountyCases)

	if min.After(cmin) {
		min = cmin
	}
	if max.Before(cmax) {
		max = cmax
	}
	return min, max
}

func WriteDateRange(db *gorm.DB) {
    min, max := dateRange(db)
    cur := min
    dates := make([]Date, 0)
    for !cur.After(max) {
        dates = append(dates, Date{cur})
        cur = cur.AddDate(0, 0, 1)
    }
    data, err := json.Marshal(dates)
    if err != nil { panic(err.Error()) }

	writeFile(string(data), "dateRange.json")
}

func str(item interface{}) string {
	bs, _ := json.Marshal(item)
	return string(bs)
}

func getDate(item interface{}) Date {
	val := reflect.ValueOf(item)
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name == "Date" {
			return val.Field(i).Interface().(Date)
		}
	}
	panic("Date not found for " + val.String())
}

func cloneWithNewDate(item interface{}, newDate time.Time) reflect.Value {
	val := reflect.ValueOf(item)
	result := reflect.New(val.Type()).Elem()
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name == "Date" {
			result.Field(i).Set(reflect.ValueOf(newDate))
		} else {
			result.Field(i).Set(val.Field(i))
		}
	}
	return result
}

func ensureRange(items interface{}, min time.Time, max time.Time) bool {
	list := reflect.ValueOf(items)
	first := list.Index(0).Interface()
	minInList := getDate(first)
	if min.Before(minInList.Time) {
		if list.Len() == 1 {
			fmt.Println("Found what looks like a brand new case for " + str(first))
			return false
		}
		panic("Missing first date for " + str(first))
	}
	last := list.Index(list.Len() - 1).Interface()
	maxInList := getDate(last)
	if max.After(maxInList.Time) {
		panic("Missing last date for " + str(last))
	}
	return true
}

func WriteWorldData(db *gorm.DB) {
	type shape struct {
		Date        time.Time `json:"date"`
		Country     string    `json:"country"`
		CountryCode string    `json:"countryCode"`
		Confirmed   int       `json:"confirmed"`
		Deaths      int       `json:"deaths"`
	}
	type mapValue struct {
		Country     string    `json:"country"`
		CountryCode string    `json:"countryCode"`
        Dates       []DateData `json:"dates"`
	}
	usBase := db.Model(&CountyCases).Where("date = data_points.date")
	usConfirmed := usBase.Select("sum(confirmed)")
	usDeaths := usBase.Select("sum(deaths)")

	countryAggregates := db.Select(`
        date,
        country,
        country_code,
        CASE
          WHEN country != 'United States' THEN sum(confirmed)
          ELSE (?)
        END as confirmed,
        CASE
          WHEN country != 'United States' THEN sum(deaths)
          ELSE (?)
        END as deaths
    `, usConfirmed.QueryExpr(), usDeaths.QueryExpr()).Model(&Point).
		Group("date, country, country_code").
		Order("date, country")

	var aggregates []shape
	countryAggregates.Scan(&aggregates)

    obj := make(map[string]mapValue)

	for _, item := range aggregates {
		val, ok := obj[item.CountryCode]
		if !ok {
            val = mapValue{
                Country: item.Country,
                CountryCode: item.CountryCode,
                Dates: make([]DateData, 0),
            }
		}
        val.Dates = append(val.Dates, makeDateData(
            item.Date, item.Confirmed, item.Deaths,
        ))
		obj[item.CountryCode] = val
	}

	min, max := dateRange(db)
	for k, val := range obj {
		if !ensureRange(val.Dates, min, max) {
			delete(obj, k)
		}
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}

	writeFile(string(bytes), "world.json")
}

func WriteProvinceData(db *gorm.DB) {
	type shape struct {
		Date        time.Time `json:"date"`
		Country     string    `json:"country"`
		CountryCode string    `json:"countryCode"`
		Province    string    `json:"province"`
		Confirmed   int       `json:"confirmed"`
		Deaths      int       `json:"deaths"`
	}
	type mapValue struct {
		Country     string    `json:"country"`
		CountryCode string    `json:"countryCode"`
		Province    string    `json:"province"`
        Dates       []DateData `json:"dates"`
	}

	query := db.Select(`
        date,
        country,
        country_code,
        province,
        sum(confirmed) as confirmed,
        sum(deaths) as deaths
    `).
		Model(&Point).
		Where("province != ''").
		Group("date, country, country_code, province").
		Order("date, country, province")

	var results []shape
	err := query.Scan(&results).Error
	if err != nil {
		panic(err.Error())
	}

	obj := make(map[string]mapValue)

	for _, item := range results {
		key := item.CountryCode + "-" + item.Province
		val, ok := obj[key]
		if !ok {
            val = mapValue{
                Country: item.Country,
                CountryCode: item.CountryCode,
                Province: item.Province,
                Dates: make([]DateData, 0),
            }
		}
        val.Dates = append(val.Dates, makeDateData(
            item.Date, item.Confirmed, item.Deaths,
        ))
		obj[key] = val
	}

	min, max := dateRange(db)
	for k, val := range obj {
		if !ensureRange(val.Dates, min, max) {
			delete(obj, k)
		}
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}

	writeFile(string(bytes), "province.json")
}

type StateFips struct {
	State  string
	FipsId string
}

func stateFipsMap(db *gorm.DB) map[string]StateFips {
	var fipsData []StateFips
	db.
		Select("state, min(substring(fips_id from 1 for 2)) as fips_id").
		Where("fips_id != ''").
		Model(&CountyCases).
		Group("state").
		Scan(&fipsData)
	fipsMap := make(map[string]StateFips)

	for _, item := range fipsData {
		fipsMap[item.State] = item
	}
	return fipsMap
}

func WriteStateData(db *gorm.DB) {
	type shape struct {
		State     string    `json:"state"`
		FipsId    string    `json:"fipsId"`
		Date      time.Time `json:"date"`
		Confirmed int       `json:"confirmed"`
		Deaths    int       `json:"deaths"`
	}
	type mapValue struct {
		State     string    `json:"state"`
		FipsId    string    `json:"fipsId"`
        Dates       []DateData `json:"dates"`
	}

	fipsMap := stateFipsMap(db)

	var results []shape
	query := db.
		Select(`
            date,
            state,
            sum(confirmed) as confirmed,
            sum(deaths) as deaths
            `).
		Where("state != ''").
		Model(&CountyCases).
		Group("date, state").
		Order("date, state")

	query.Scan(&results)

	obj := make(map[string]mapValue)

	for _, item := range results {
		stateFips := fipsMap[item.State]
		if stateFips.FipsId == "" {
			continue
		}

		item.FipsId = stateFips.FipsId
		val, ok := obj[item.FipsId]
		if !ok {
            val = mapValue{
                State: item.State,
                FipsId: item.FipsId,
                Dates: make([]DateData, 0),
            }
		}
        val.Dates = append(val.Dates, makeDateData(
            item.Date, item.Confirmed, item.Deaths,
        ))
		obj[item.FipsId] = val
	}

	min, max := dateRange(db)
	for k, val := range obj {
		if !ensureRange(val.Dates, min, max) {
			delete(obj, k)
		}
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}

	writeFile(string(bytes), "state.json")
}

func WriteCountyData(db *gorm.DB) {
	type shape struct {
		State  string `json:"state"`
		County string `json:"county"`
		FipsId string `json:"fipsId"`
		Date      time.Time `json:"date"`
		Confirmed int       `json:"confirmed"`
		Deaths    int       `json:"deaths"`
	}
	type mapValue struct {
        Id string `json:"id"`
        Dates       []DateData `json:"dates"`
	}
	query := db.
		Select(`
            date,
            fips_id,
            state,
            county,
            sum(confirmed) as confirmed,
            sum(deaths) as deaths
        `).
		Model(&CountyCases).
		Where("fips_id != ''").
		Group("date, state, county, fips_id").
		Order("date, state, county, fips_id")

	var results []shape
	query.Scan(&results)

    obj := make(map[string]mapValue)

	for _, item := range results {
        val, ok := obj[item.FipsId]
		if !ok {
            val = mapValue{
                Id: item.FipsId,
                Dates: make([]DateData, 0),
            }
		}
        val.Dates = append(val.Dates, makeDateData(
            item.Date, item.Confirmed, item.Deaths,
        ))
		obj[item.FipsId] = val
	}


	min, max := dateRange(db)
	for k, val := range obj {
		if !ensureRange(val.Dates, min, max) {
			delete(obj, k)
		}
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	writeFile(string(bytes), "county.json")
}

