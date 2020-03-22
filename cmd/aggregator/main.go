package main

import (
	"covid-tracker/data"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Aggregate struct {
    Country string
    Confirmed int
}

func main() {
    db, err := gorm.Open("sqlite3", "test.db")
    db.LogMode(true)
    if err != nil {
		fmt.Println(err)
        panic("failed to connect database")
    }

    aggregates := make([]Aggregate, 0)
    db.Model(&data.Point).
        Select("sum(confirmed) as confirmed, country").
        Where("date(date) = (?)", db.Model(&data.Point).Select("max(date(date))").QueryExpr()).
        Group("country").
        Order("confirmed desc").
        Scan(&aggregates)

    for _, aggregate := range aggregates {
        fmt.Printf("country: %s, confirmed: %d\n", aggregate.Country, aggregate.Confirmed)
    }
}
