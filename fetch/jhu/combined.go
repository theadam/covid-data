package jhu

import (
	"covid-tracker/data"
	"fmt"
	"time"
)

func maxDate(timeData []*data.DataPoint) time.Time {
	var max time.Time
	for _, data := range timeData {
        if max.Before(data.Date) {
			max = data.Date
		}
	}
	return max
}

// Get the next date to use after timeseries data
func nextDate(timeData []*data.DataPoint) time.Time {
    max := maxDate(timeData)
    // 1 day after latest data
    next := max.AddDate(0, 0, 1)
    if time.Until(next).Hours() > 24 {
        fmt.Println(next.String())
        panic("Next date is too far in the future")
    }
    return next
}

func collectLatest(values []*data.DataPoint) map[data.LocationKey]*data.DataPoint {
    latest := make(map[data.LocationKey]*data.DataPoint)
    for _, value := range values {
        key := *value.LocationKey()
        current, ok := latest[key]
        if !ok || current.Date.Before(value.Date) {
            latest[key] = value
        }
    }
    return latest
}

func GetData() []*data.DataPoint {
    globals := timeSeriesGlobals()
    us := timeSeriesUs()

    combined := append(globals, us...)
    next := nextDate(combined)

    currentData := fetchCurrent(&next, collectLatest(combined))

	return append(combined, currentData...)
}
