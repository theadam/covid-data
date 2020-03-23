package main

import (
	"covid-tracker/data"
	"covid-tracker/fetch/jhu"
	"covid-tracker/fetch/opta"
	"covid-tracker/utils"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
)

func loadJhu(db *gorm.DB) {
	points, err := jhu.GetData()
	if err != nil {
		fmt.Println(err)
		return
	}

    db.Transaction(func(tx *gorm.DB) error {
        db.Unscoped().Delete(&data.Point)
        for _, point := range points {
            tx.Create(&point)
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

    db.Transaction(func(tx *gorm.DB) error {
        db.Unscoped().Delete(&data.CountyCases)
        for _, item := range countyDatas {
            tx.Create(&item)
        }
        return nil
    })
}

func main() {
    db := utils.OpenDB()
    defer db.Close()

    db.AutoMigrate(&data.Point, &data.CountyCases)


    skipJhu := flag.Bool("skip-jhu", false, "Skip loading johns hopkins university data")
    skipOpta := flag.Bool("skip-1point3acres", false, "Skip 1point3acres data")

    flag.Parse()

    if !*skipJhu {
        loadJhu(db)
    }
    if !*skipOpta {
        loadOpta(db)
    }
}
