package main

import (
	"covid-tracker/fetch"
	"covid-tracker/data"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Aggregate struct {
    country string
    confirmed int
}

func main() {
    db, err := gorm.Open("sqlite3", ":memory:")
    db.LogMode(true)
    if err != nil {
		fmt.Println(err)
        panic("failed to connect database")
    }

	points, err := fetch.GetData()
	if err != nil {
		fmt.Println(err)
		return
	}

    db.AutoMigrate(&data.Point)

    for _, point := range points {
        fmt.Printf("State: %s\n", point.State)
        fmt.Printf("Country: %s\n", point.Country)
        fmt.Printf("Confirmed: %d\n", point.Confirmed)
        fmt.Printf("Deaths: %d\n", point.Deaths)
        fmt.Printf("Recovered: %d\n", point.Recovered)
        fmt.Printf("Date: %s\n", point.Date)
        fmt.Printf("Lat: %s\n", point.Lat)
        fmt.Printf("Long: %s\n", point.Long)
        fmt.Println()

        db.Create(&point)
    }

    aggregates := make([]Aggregate, 0)
    db.Model(&data.Point).
        Select("sum(confirmed) as confirmed, country").
        Where("date(date) = ?", db.Model(&data.Point).Select("max(date(date, '-1 day'))").QueryExpr()).
        Group("country").
        Scan(&aggregates)
    fmt.Println(aggregates)
}
