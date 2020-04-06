package data

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

func LoadWorldTable(db *gorm.DB) {
	type shape struct {
		Date      string `json:"date"`
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
		slice, ok := obj[item.Country]
		if !ok {
			slice = make([]shape, 0)
		}
		slice = append(slice, item)
		obj[item.Country] = slice
	}

    bytes, err := json.Marshal(obj)
    if err != nil { panic(err.Error()) }

    err = db.Create(&WorldHistorical{ Data: string(bytes) }).Error
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
		Date      string `json:"date"`
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
		Date      string `json:"date"`
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

    bytes, err := json.Marshal(obj)
    if err != nil { panic(err.Error()) }

    err = db.Create(&CountyHistorical{ Data: string(bytes) }).Error
    if err != nil { panic(err.Error()) }
}
