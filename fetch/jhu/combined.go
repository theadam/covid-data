package jhu

import (
	"covid-tracker/data"
	"fmt"
	"time"
)

func maxDate(timeData []TimeValue) time.Time {
	var max time.Time
	for _, data := range timeData {
        if max.Before(data.Date) {
			max = data.Date
		}
	}
	return max
}

// Get the next date to use after timeseries data
func nextDate(timeData []TimeValue) time.Time {
    max := maxDate(timeData)
    // 1 day after latest data
    next := max.AddDate(0, 0, 1)
    if next.Sub(time.Now()).Hours() > 24 {
        fmt.Println(next.String())
        panic("Next date is too far in the future")
    }
    return next
}

func GetData(start time.Time) ([]data.DataPoint, []data.CountyData) {
    globals := timeSeriesGlobals(start)
    nextGlobalDate := nextDate(globals)

    us := timeSeriesUs(start)
    nextUsDate := nextDate(us)

    globalData, usData := fetchCurrent(nextGlobalDate, nextUsDate)

    for _, global := range globals {
        globalData = append(globalData, toDataPoint(global));
    }

    for _, county := range us {
        usData = append(usData, toCountyData(county));
    }

	return globalData, usData
}
