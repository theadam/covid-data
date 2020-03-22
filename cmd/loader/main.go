package main

import (
	"covid-tracker/fetch"
	"covid-tracker/data"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
    db, err := gorm.Open("sqlite3", "test.db")
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

    db.DropTableIfExists(&data.Point)
    db.AutoMigrate(&data.Point)

    db.Transaction(func(tx *gorm.DB) error {
        for _, point := range points {
            tx.Create(&point)
        }
        return nil
    })
}
