package data

import (
    "github.com/jinzhu/gorm"
	"time"
)

type DataPoint struct {
    gorm.Model
    State string
    Country string
    Confirmed int
    Deaths int
    Recovered int
    Date time.Time
    Lat string
    Long string
}

var Point DataPoint


