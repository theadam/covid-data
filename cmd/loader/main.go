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

func loadJhu(db *gorm.DB, ignoreStart bool) {
    var start time.Time
    if !ignoreStart { start = startDate(db, &data.Point) }

	points, counties := jhu.GetData(start)

    // INSERT GLOBAL DATA
    items := make([]interface{}, len(points))
    for i, v := range points {
        items[i] = v
    }
    fmt.Println("Inserting " + strconv.Itoa(len(items)) + " items")

    db.Unscoped().Where("date >= ?", start).Delete(&data.Point)
    err := gormbulk.BulkInsert(db, items, 1000)
    if err != nil { panic(err.Error()) }


    // INSERT US DATA
    if !ignoreStart { start = startDate(db, &data.CountyCases) }
    items = make([]interface{}, len(counties))
    for i, v := range counties {
        items[i] = v
    }
    fmt.Println("Inserting " + strconv.Itoa(len(items)) + " items")

    db.Unscoped().Where("date >= ?", start).Delete(&data.CountyCases)
    err = gormbulk.BulkInsert(db, items, 1000)
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

    db.AutoMigrate(&data.Point, &data.CountyCases)

    ignoreStart := flag.Bool("all-dates", false, "Ignore start date")
    flag.Parse()

    db.Transaction(func(tx *gorm.DB) error {
        loadJhu(tx, *ignoreStart)
        return nil
    })
}
