package main

import (
	"covid-tracker/data"
	"covid-tracker/fetch/jhu"
	"fmt"
	"time"
)

func runAction(name string, action func()) {
	now := time.Now()
	fmt.Println("Starting " + name)
	action()
	fmt.Println("Finished " + name + " in " + time.Since(now).String())
	fmt.Println()
}

func main() {
	fmt.Println("Loading all data")
	fmt.Println()

	// Gets all data for all time
	var start time.Time

	points, counties := jhu.GetData(start)

	runAction("Writing World JSON data", func() { data.WriteWorldData(points, counties) })
	runAction("Writing Province JSON data", func() { data.WriteProvinceData(points, counties) })
	runAction("Writing State JSON data", func() { data.WriteStateData(points, counties) })
	runAction("Writing County JSON data", func() { data.WriteCountyData(points, counties) })
	runAction("Writing Date Range JSON data", func() { data.WriteDateRange(points, counties) })

}
