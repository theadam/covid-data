package server

import (
	"covid-tracker/data"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Env struct {
	db *gorm.DB
}

func (env *Env) GetCountries(c *gin.Context) {
	type shape struct {
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
	}
	var countries []shape
	env.db.Model(&data.Point).Select("country, country_code").Group("country, country_code").Scan(&countries)

	c.JSON(200, countries)
}

func (env *Env) GetStates(c *gin.Context) {
	var states []string
	env.db.Model(&data.CountyCases).Select("state").Group("state").Pluck("state", &states)

	c.JSON(200, states)
}

func (env *Env) GetCounties(c *gin.Context) {
	type shape struct {
		state  string
		county string
	}

	var result []shape
	env.db.Model(&data.CountyCases).Select("state, county").Group("state, county").Order("state, county").Scan(&result)

	c.JSON(200, result)
}

func (env *Env) GetCountryData(c *gin.Context) {
	type shape struct {
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		Confirmed   int    `json:"confirmed"`
		Deaths      int    `json:"deaths"`
	}
	var results []shape

	maxDate := env.db.Model(&data.Point).Select("max(date)").QueryExpr()

	env.db.
		Model(&data.Point).
		Select("country, country_code, sum(confirmed) as confirmed, sum(deaths) as deaths").
		Group("country, country_code").
		Where("date = (?)", maxDate).
		Scan(&results)

	var usResult shape
	env.db.Model(&data.CountyCases).
		Select("sum(confirmed) as confirmed, sum(deaths) as deaths").
		Scan(&usResult)

	for i, result := range results {
		if result.Country == "United States" {
			result.Confirmed = usResult.Confirmed
			result.Deaths = usResult.Deaths
			results[i] = result
		}
	}

	c.JSON(200, results)
}

func countryAggregates(db *gorm.DB) *gorm.DB {
	usBase := db.Model(&data.CountyCases).Where("date <= data_points.date")
	usConfirmed := usBase.Select("sum(confirmed)")
	usDeaths := usBase.Select("sum(deaths)")

	countryAggregates := db.Select(`
        date,
        country,
        country_code,
        CASE
          WHEN country != "United States" THEN sum(confirmed)
          ELSE (?)
        END as confirmed,
        CASE
          WHEN country != "United States" THEN sum(deaths)
          ELSE (?)
        END as deaths
    `, usConfirmed.QueryExpr(), usDeaths.QueryExpr()).Model(&data.Point).
		Group("date, country, country_code").
		Order("date, country")

	return countryAggregates
}

func (env *Env) GetCountryHistorical(c *gin.Context) {
	type shape struct {
		Date        time.Time `json:"date"`
		Country     string    `json:"country"`
		CountryCode string    `json:"countryCode"`
		Confirmed   int       `json:"confirmed"`
		Deaths      int       `json:"deaths"`
	}

	query := countryAggregates(env.db)
	if c.Query("country") != "" {
		query = query.Where("country IN (?)", strings.Split(c.Query("country"), ","))
	}
	var aggregates []shape
	query.Scan(&aggregates)

	obj := make(map[string][]shape)

	for _, item := range aggregates {
		slice, ok := obj[item.Country]
		if !ok {
			slice = make([]shape, 0)
		}
		slice = append(slice, item)
		obj[item.Country] = slice
	}

	c.JSON(200, obj)
}

func (env *Env) GetStateData(c *gin.Context) {
	type shape struct {
		State     string `json:"state"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
	}

	var results []shape
	env.db.
		Select(`
            CASE WHEN state = "" THEN "Unknown" ELSE state END as state, sum(confirmed) as confirmed, sum(deaths) as deaths
        `).
		Model(&data.CountyCases).
		Group("state").
		Scan(&results)

	c.JSON(200, results)
}

func (env *Env) GetStateHistorical(c *gin.Context) {
	type shape struct {
		State     string `json:"state"`
		Date      string `json:"date"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
	}

	base := env.db.Model(&data.CountyCases).Where("date <= outer_data.date").Where("state = states.state")
	confirmed := base.Select("sum(confirmed)")
	deaths := base.Select("sum(deaths)")

	tableName := env.db.NewScope(&data.CountyCases).TableName()

	states := env.db.Model(&data.CountyCases).Select("state").Group("state")

	var results []shape
	query := env.db.
		Select(`
            date,
            CASE WHEN states.state = "" THEN "Unknown" ELSE states.state END as state,
            (?) as confirmed,
            (?) as deaths
        `, confirmed.QueryExpr(), deaths.QueryExpr()).
		Table(tableName+" outer_data").
		Joins("CROSS JOIN (?) states", states.QueryExpr()).
		Group("date, states.state").
		Order("date, states.state")

	if c.Query("state") != "" {
		query = query.Where("states.state in (?)", strings.Split(c.Query("state"), ","))
	}

	query.Scan(&results)

	c.JSON(200, results)
}

func (env *Env) GetCountyData(c *gin.Context) {
	type shape struct {
		State     string `json:"state"`
		County    string `json:"county"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
	}

	var results []shape
	env.db.
		Select(`
            CASE WHEN county = "" THEN "Unknown" ELSE county END as county,
            CASE WHEN state = "" THEN "Unknown" ELSE state END as state,
            sum(confirmed) as confirmed,
            sum(deaths) as deaths
        `).
		Model(&data.CountyCases).
		Group("state, county").
		Scan(&results)

	c.JSON(200, results)
}

func (env *Env) GetCountyHistorical(c *gin.Context) {
	type shape struct {
		State     string `json:"state"`
		County    string `json:"county"`
		Date      string `json:"date"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
	}

	base := env.db.Model(&data.CountyCases).Where("date <= outer_data.date").Where("state = counties.state").Where("county = counties.county")
	confirmed := base.Select("sum(confirmed)")
	deaths := base.Select("sum(deaths)")

	tableName := env.db.NewScope(&data.CountyCases).TableName()

	counties := env.db.Model(&data.CountyCases).Select("state, county").Group("state, county")

	var results []shape
	query := env.db.
		Select(`
            date,
            CASE WHEN counties.state = "" THEN "Unknown" ELSE counties.state END as state,
            CASE WHEN counties.county = "" THEN "Unknown" ELSE counties.county END as county,
            (?) as confirmed,
            (?) as deaths
        `, confirmed.QueryExpr(), deaths.QueryExpr()).
		Table(tableName+" outer_data").
		Joins("CROSS JOIN (?) counties", counties.QueryExpr()).
		Group("date, counties.state, counties.county").
		Order("date, counties.state, counties.county")

	if c.Query("state") != "" {
		query = query.Where("counties.state in (?)", strings.Split(c.Query("state"), ","))
	}
	if c.Query("county") != "" {
		query = query.Where("counties.county in (?)", strings.Split(c.Query("county"), ","))
	}

	query.Scan(&results)

	c.JSON(200, results)
}

func Router(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	env := &Env{db: db}

	r.GET("/countries", env.GetCountries)
	r.GET("/states", env.GetStates)
	r.GET("/counties", env.GetCounties)
	r.GET("/data/countries", env.GetCountryData)
	r.GET("/data/countries/historical", env.GetCountryHistorical)
	r.GET("/data/us/states", env.GetStateData)
	r.GET("/data/us/states/historical", env.GetStateHistorical)
	r.GET("/data/us/counties", env.GetCountyData)
	r.GET("/data/us/counties/historical", env.GetCountyHistorical)
	return r
}
