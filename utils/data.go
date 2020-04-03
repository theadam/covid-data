package utils

import (
	"covid-tracker/data"
	"time"
)

func FilterDataPointByDate(points []data.DataPoint, start time.Time) []data.DataPoint {
    result := make([]data.DataPoint, 0)
    for _, p := range points {
        if p.Date.After(start) || p.Date.Equal(start) {
            result = append(result, p)
        }
    }
    return result
}

func FilterCountyDataByDate(points []data.CountyData, start time.Time) []data.CountyData {
    result := make([]data.CountyData, 0)
    for _, p := range points {
        if p.Date.After(start) || p.Date.Equal(start) {
            result = append(result, p)
        }
    }
    return result
}
