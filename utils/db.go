package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

const waittime = 40

func OpenDB() *gorm.DB {
    cur := 0
    var db *gorm.DB
    var err error
    for (err != nil || db == nil) && cur < waittime {
        db, err := gorm.Open("postgres", "host=postgres user=postgres dbname=covid password=password sslmode=disable")
        if err == nil {
            db.LogMode(true)
            return db
        }
    }
    fmt.Println(err)
    panic("failed to connect database")
}
