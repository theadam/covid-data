package main

import (
	"covid-tracker/fetch/jhu"
	"covid-tracker/fetch/opta"
	"covid-tracker/data"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func loadJhu(db *gorm.DB) {
	points, err := jhu.GetData()
	if err != nil {
		fmt.Println(err)
		return
	}

    db.Transaction(func(tx *gorm.DB) error {
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
    fmt.Print(len(countyDatas))
}

func main() {
    db, err := gorm.Open("sqlite3", "test.db")
    db.LogMode(true)
    if err != nil {
		fmt.Println(err)
        panic("failed to connect database")
    }
    db.AutoMigrate(&data.Point, &data.CountyCases)

    // loadJhu()
    loadOpta(db)
}
