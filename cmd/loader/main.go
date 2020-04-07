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

	db.AutoMigrate(&data.Point, &data.CountyCases, &data.WorldHist, &data.CountyHist, &data.StateHist, &data.ProvinceHist)

	ignoreStart := flag.Bool("all-dates", false, "Ignore start date")
	flag.Parse()

	if *ignoreStart {
		fmt.Println("Loading all data")
	} else {
		fmt.Println("Reloading recent data")
	}
	fmt.Println()

	var globalStart time.Time
	if !*ignoreStart {
		globalStart = startDate(db, &data.Point)
	}

	var usStart time.Time
	if !*ignoreStart {
		usStart = startDate(db, &data.CountyCases)
	}

	points, counties := jhu.GetData(globalStart, usStart)

	db.Transaction(func(tx *gorm.DB) error {
		runAction("Loading Globals", func() { loadGlobals(tx, points, globalStart) })
		runAction("Loading US", func() { loadUs(tx, counties, usStart) })
		runAction("Loading World Cache Table", func() { data.LoadWorldTable(tx) })
		runAction("Loading State Cache Table", func() { data.LoadStateTable(tx) })
		runAction("Loading County Cache Table", func() { data.LoadCountyTable(tx) })
		runAction("Loading Province Cache Table", func() { data.LoadProvinceTable(tx) })

		return nil
	})
}
