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
        State  string `json:"state"`
        County string `json:"county"`
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
	maxCountyDate := env.db.Model(&data.CountyCases).Select("max(date)").QueryExpr()

	env.db.
		Model(&data.Point).
		Select("country, country_code, sum(confirmed) as confirmed, sum(deaths) as deaths").
		Group("country, country_code").
		Where("date = (?)", maxDate).
		Scan(&results)

	var usResult shape
	env.db.Model(&data.CountyCases).
		Select("sum(confirmed) as confirmed, sum(deaths) as deaths").
        Where("date = (?)", maxCountyDate).
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
	usBase := db.Model(&data.CountyCases).Where("date = data_points.date")
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
	maxDate := env.db.Model(&data.CountyCases).Select("max(date)").QueryExpr()

	var results []shape
	env.db.
		Select(`
            CASE WHEN state = '' THEN 'Unknown' ELSE state END as state, sum(confirmed) as confirmed, sum(deaths) as deaths
        `).
        Where("date = (?)", maxDate).
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

	var results []shape
	query := env.db.
		Select(`
            date,
            CASE WHEN state = '' THEN 'Unknown' ELSE state END as state,
            sum(confirmed) as confirmed,
            sum(deaths) as deaths
            `).
            Model(&data.CountyCases).
            Group("date, state").
            Order("date, state")

	if c.Query("state") != "" {
		query = query.Where("state in (?)", strings.Split(c.Query("state"), ","))
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

	maxDate := env.db.Model(&data.CountyCases).Select("max(date)").QueryExpr()
	var results []shape
	env.db.
		Select(`
            CASE WHEN county = '' THEN 'Unknown' ELSE county END as county,
            CASE WHEN state = '' THEN 'Unknown' ELSE state END as state,
            sum(confirmed) as confirmed,
            sum(deaths) as deaths
        `).
		Model(&data.CountyCases).
        Where("date = (?)", maxDate).
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
		FipsId    string    `json:"fipsId"`
	}

	var results []shape
	query := env.db.
		Select(`
            date,
            fips_id,
            state,
            county,
            sum(confirmed) as confirmed,
            sum(deaths) as deaths
        `).
        Model(&data.CountyCases).
        Where("fips_id != ''").
		Group("date, state, county, fips_id").
		Order("date, state, county, fips_id")

	if c.Query("state") != "" {
		query = query.Where("state in (?)", strings.Split(c.Query("state"), ","))
	}
	if c.Query("county") != "" {
		query = query.Where("county in (?)", strings.Split(c.Query("county"), ","))
	}

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


	c.JSON(200, obj)
}

func Router(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	env := &Env{db: db}

    api := r.Group("/api/")
    {
        api.GET("/countries", env.GetCountries)
        api.GET("/states", env.GetStates)
        api.GET("/counties", env.GetCounties)
        api.GET("/data/countries", env.GetCountryData)
        api.GET("/data/countries/historical", env.GetCountryHistorical)
        api.GET("/data/us/states", env.GetStateData)
        api.GET("/data/us/states/historical", env.GetStateHistorical)
        api.GET("/data/us/counties", env.GetCountyData)
        api.GET("/data/us/counties/historical", env.GetCountyHistorical)
    }

    r.Static("/client", "./client/build")

    r.NoRoute(func(c *gin.Context) {
        // c.String(404, "404 page not found")

        c.File("./client/build/index.html")
    })

	return r
}

