package main

import (
	"covid-tracker/data"
	"covid-tracker/fetch/jhu"
	"covid-tracker/fetch/opta"
	"covid-tracker/utils"
	"flag"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/t-tiger/gorm-bulk-insert"
)

func loadJhu(db *gorm.DB) {
	points, err := jhu.GetData()
	if err != nil {
		fmt.Println(err)
		return
	}

    items := make([]interface{}, len(points))
    for i, v := range points {
        items[i] = v
    }
    fmt.Println("Inserting " + strconv.Itoa(len(items)) + " items")

    db.Transaction(func(tx *gorm.DB) error {
        tx.Unscoped().Delete(&data.Point)
        err := gormbulk.BulkInsert(tx, items, 1000)
        if err != nil {
            panic(err.Error())
        }

        return nil
    })
}

func loadOpta(db *gorm.DB) {
    countyDatas, err := opta.GetData()
	if err != nil {
		fmt.Println(err)
		return
	}

    items := make([]interface{}, len(countyDatas))
    for i, v := range countyDatas {
        items[i] = v
    }
    fmt.Println("Inserting " + strconv.Itoa(len(items)) + " items")

    db.Transaction(func(tx *gorm.DB) error {
        tx.Unscoped().Delete(&data.CountyCases)
        err := gormbulk.BulkInsert(tx, items, 1000)
        if err != nil {
            panic(err.Error())
        }

        return nil
    })
}

func main() {
    db := utils.OpenDB()
    defer db.Close()

    db.AutoMigrate(&data.Point, &data.CountyCases)


    runJhu := flag.Bool("jhu", false, "Load johns hopkins university data")
    runOpta := flag.Bool("opta", false, "Load 1point3acres data")

    flag.Parse()

    runAll := !*runJhu && !*runOpta

    if runAll || *runJhu {
        loadJhu(db)
    }
    if runAll || *runOpta {
        loadOpta(db)
    }
}
