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

	points := jhu.GetData()

	runAction("Writing World JSON data", func() { data.WriteWorldData(points) })
	runAction("Writing Province JSON data", func() { data.WriteProvinceData(points) })
	runAction("Writing State JSON data", func() { data.WriteStateData(points) })
	runAction("Writing County JSON data", func() { data.WriteCountyData(points) })
	runAction("Writing Date Range JSON data", func() { data.WriteDateRange(points) })
}
