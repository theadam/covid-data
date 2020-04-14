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
	fmt.Println("Inserting " + strconv.Itoa(len(items)) + " global items")

	db.Unscoped().Where("date >= ?", start).Delete(&data.Point)
	err := gormbulk.BulkInsert(db, items, 3000)
	if err != nil {
		panic(err.Error())
	}
}

func loadUs(db *gorm.DB, counties []data.CountyData, start time.Time) {
	// INSERT US DATA
	items := make([]interface{}, len(counties))
	for i, v := range counties {
		items[i] = v
	}
	fmt.Println("Inserting " + strconv.Itoa(len(items)) + " US items")

	db.Unscoped().Where("date >= ?", start).Delete(&data.CountyCases)
	err := gormbulk.BulkInsert(db, items, 3000)
	if err != nil {
		panic(err.Error())
	}
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

func runAction(name string, action func()) {
	now := time.Now()
	fmt.Println("Starting " + name)
	action()
	fmt.Println("Finished " + name + " in " + time.Since(now).String())
	fmt.Println()
}

func main() {
	db := utils.OpenDB()
	defer db.Close()

	db.AutoMigrate(&data.Point, &data.CountyCases)

	flag.Parse()

    fmt.Println("Loading all data")
	fmt.Println()

    // Gets all data for all time
	var start time.Time

	points, counties := jhu.GetData(start)

	db.Transaction(func(tx *gorm.DB) error {
		runAction("Loading Globals", func() { loadGlobals(tx, points, start) })
		runAction("Loading US", func() { loadUs(tx, counties, start) })
		runAction("Writing World JSON data", func() { data.WriteWorldData(tx) })
		runAction("Writing Province JSON data", func() { data.WriteProvinceData(tx) })
		runAction("Writing State JSON data", func() { data.WriteStateData(tx) })
		runAction("Writing County JSON data", func() { data.WriteCountyData(tx) })
		runAction("Writing Date Range JSON data", func() { data.WriteDateRange(tx) })

		return nil
	})
}
