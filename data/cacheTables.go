package data

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

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

func str(item interface{}) string {
    bs, _ := json.Marshal(item)
    return string(bs)
}

func getDate(item interface{}) time.Time {
    val := reflect.ValueOf(item)
    for i := 0; i < val.NumField(); i ++ {
        if (val.Type().Field(i).Name == "Date") {
            return val.Field(i).Interface().(time.Time)
        }
    }
    panic("Date not found for " + val.String())
}

func cloneWithNewDate(item interface{}, newDate time.Time) reflect.Value {
    val := reflect.ValueOf(item)
    result := reflect.New(val.Type()).Elem()
    for i := 0; i < val.NumField(); i ++ {
        if (val.Type().Field(i).Name == "Date") {
            result.Field(i).Set(reflect.ValueOf(newDate))
        } else {
            result.Field(i).Set(val.Field(i))
        }
    }
    return result
}

func ensureRange(items interface{}, min time.Time, max time.Time) interface{} {
    list := reflect.ValueOf(items)
    first := list.Index(0).Interface()
    minInList := getDate(first)
    if min.Before(minInList) {
        panic("Missing first date for " + str(first))
    }
    last := list.Index(list.Len() - 1).Interface()
    maxInList := getDate(last)
    if max.After(maxInList) {
        panic("Missing last date for " + str(last))
    }
    return list.Interface()
}


func LoadWorldTable(db *gorm.DB) {
	type shape struct {
		Date      time.Time `json:"date"`
		Country    string `json:"country"`
		CountryCode    string `json:"countryCode"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
	}
    db.Unscoped().Delete(&WorldHist)
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

	obj := make(map[string][]shape)

	for _, item := range aggregates {
		slice, ok := obj[item.CountryCode]
		if !ok {
			slice = make([]shape, 0)
		}
		slice = append(slice, item)
		obj[item.CountryCode] = slice
	}


    min, max := dateRange(db)
    for _, val := range obj {
        ensureRange(val, min, max)
    }

    bytes, err := json.Marshal(obj)
    if err != nil { panic(err.Error()) }

    err = db.Create(&WorldHistorical{ Data: string(bytes) }).Error
    if err != nil { panic(err.Error()) }
}

func LoadProvinceTable(db *gorm.DB) {
	type shape struct {
		Date      time.Time `json:"date"`
		Country    string `json:"country"`
		CountryCode    string `json:"countryCode"`
		Province    string `json:"province"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
	}
    db.Unscoped().Delete(&ProvinceHist)

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
    if err != nil { panic(err.Error()) }

	obj := make(map[string][]shape)

	for _, item := range results {
        key := item.CountryCode + "-" + item.Province
		slice, ok := obj[key]
		if !ok {
			slice = make([]shape, 0)
		}
		slice = append(slice, item)
		obj[key] = slice
	}

    min, max := dateRange(db)
    for _, val := range obj {
        ensureRange(val, min, max)
    }

    bytes, err := json.Marshal(obj)
    if err != nil { panic(err.Error()) }

    err = db.Create(&ProvinceHistorical{ Data: string(bytes) }).Error
    if err != nil { panic(err.Error()) }
}

type StateFips struct {
    State     string
    FipsId    string
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

func LoadStateTable(db *gorm.DB) {
    db.Unscoped().Delete(&StateHist)
	type shape struct {
		State     string `json:"state"`
		FipsId    string    `json:"fipsId"`
		Date      time.Time `json:"date"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
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


	obj := make(map[string][]shape)

	for _, item := range results {
        stateFips := fipsMap[item.State]
        if stateFips.FipsId == "" { continue }

        item.FipsId = stateFips.FipsId
		slice, ok := obj[item.FipsId]
		if !ok {
			slice = make([]shape, 0)
		}
		slice = append(slice, item)
		obj[item.FipsId] = slice
	}

    min, max := dateRange(db)
    for _, val := range obj {
        ensureRange(val, min, max)
    }

    bytes, err := json.Marshal(obj)
    if err != nil { panic(err.Error()) }

    err = db.Create(&StateHistorical{ Data: string(bytes) }).Error
    if err != nil { panic(err.Error()) }
}

func LoadCountyTable(db *gorm.DB) {
    db.Unscoped().Delete(&CountyHist)

	type shape struct {
		State     string `json:"state"`
		County    string `json:"county"`
		Date      time.Time `json:"date"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
		FipsId    string    `json:"fipsId"`
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

	obj := make(map[string][]shape)

	for _, item := range results {
		slice, ok := obj[item.FipsId]
		if !ok {
			slice = make([]shape, 0)
		}
		slice = append(slice, item)
		obj[item.FipsId] = slice
	}

    min, max := dateRange(db)
    for _, val := range obj {
        ensureRange(val, min, max)
    }

    bytes, err := json.Marshal(obj)
    if err != nil { panic(err.Error()) }

    err = db.Create(&CountyHistorical{ Data: string(bytes) }).Error
    if err != nil { panic(err.Error()) }
}
