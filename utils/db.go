package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const waittime = 10

func OpenDB() *gorm.DB {
    cur := 0
    var db *gorm.DB
    var err error

    for (err != nil || db == nil) && cur < waittime {
        db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL") + "?sslmode=disable")
        if err == nil {
            db.LogMode(true)
            return db
        }
        fmt.Println("Failed to connect to database... Trying again.")
        time.Sleep(1 * time.Second)
        cur++
    }
    fmt.Println(err.Error())
    panic("failed to connect database")
}
