package main

import (
	"covid-tracker/data"
	"covid-tracker/fetch/jhu"
	"covid-tracker/utils"
	"flag"
	"fmt"
	"strconv"
    "time"

	"github.com/jinzhu/gorm"
	"github.com/t-tiger/gorm-bulk-insert"
)

func loadGlobals(db *gorm.DB, globals []data.DataPoint, start time.Time) {
    // INSERT GLOBAL DATA
    items := make([]interface{}, len(globals))
    for i, v := range globals {
        items[i] = v
    }
    fmt.Println("Inserting " + strconv.Itoa(len(items)) + " items")

    db.Unscoped().Where("date >= ?", start).Delete(&data.Point)
    err := gormbulk.BulkInsert(db, items, 1000)
    if err != nil { panic(err.Error()) }
}

func loadUs(db *gorm.DB, counties []data.CountyData, start time.Time) {
    // INSERT US DATA
    items := make([]interface{}, len(counties))
    for i, v := range counties {
        items[i] = v
    }
    fmt.Println("Inserting " + strconv.Itoa(len(items)) + " items")

    db.Unscoped().Where("date >= ?", start).Delete(&data.CountyCases)
    err := gormbulk.BulkInsert(db, items, 1000)
    if err != nil { panic(err.Error()) }
}

func startDate(db *gorm.DB, table interface{}) time.Time {
    var dates []time.Time
    db.Model(table).Select("max(date) as date").Pluck("date", &dates)
    if len(dates) > 0 {
        return dates[0].AddDate(0, 0, -1)
    }
    var zero time.Time
    return zero
}

func main() {
    db := utils.OpenDB()
    defer db.Close()

    db.AutoMigrate(&data.Point, &data.CountyCases, &data.WorldHist, &data.CountyHist, &data.StateHist)

    ignoreStart := flag.Bool("all-dates", false, "Ignore start date")
    flag.Parse()

    var globalStart time.Time
    if !*ignoreStart { globalStart = startDate(db, &data.Point) }

    var usStart time.Time
    if !*ignoreStart { usStart = startDate(db, &data.CountyCases) }

	points, counties := jhu.GetData(globalStart, usStart)

    db.Transaction(func(tx *gorm.DB) error {
        loadGlobals(tx, points, globalStart)
        loadUs(tx, counties, usStart)

        data.LoadWorldTable(tx)
        data.LoadStateTable(tx)
        data.LoadCountyTable(tx)
        return nil
    })
}
