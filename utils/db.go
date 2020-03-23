package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func OpenDB() *gorm.DB {
    db, err := gorm.Open("sqlite3", "test.db")
    db.LogMode(true)
    if err != nil {
		fmt.Println(err)
        panic("failed to connect database")
    }
    return db
}
