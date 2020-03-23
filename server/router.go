package server

import (
	"covid-tracker/data"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Env struct {
	db *gorm.DB
}

func (env *Env) GetCountries(c *gin.Context) {
	var countries []string
	env.db.Model(&data.Point).Select("country").Group("country").Pluck("country", &countries)

	c.JSON(200, countries)
}

func (env *Env) GetStates(c *gin.Context) {
	var states []string
	env.db.Model(&data.CountyCases).Select("state").Group("state").Pluck("state", &states)

	c.JSON(200, states)
}

func (env *Env) GetCountryData(c *gin.Context) {
	type shape struct {
		Country   string `json:"country"`
		Confirmed int    `json:"confirmed"`
		Deaths    int    `json:"deaths"`
		Recovered int    `json:"recovered"`
	}
	var results []shape

    maxDate := env.db.Model(&data.Point).Select("max(date)").QueryExpr()

	env.db.
		Model(&data.Point).
		Select("country, sum(confirmed) as confirmed, sum(deaths) as deaths, sum(recovered) as recovered").
		Group("country").
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

func Router(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	env := &Env{db: db}

	r.GET("/countries", env.GetCountries)
	r.GET("/states", env.GetStates)
	r.GET("/data/country", env.GetCountryData)
	return r
}
