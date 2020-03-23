package main

import (
	"covid-tracker/data"
	"covid-tracker/utils"
	"fmt"
)

type Aggregate struct {
    Country string
    Confirmed int
}

func main() {
    db := utils.OpenDB()
    defer db.Close()

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
