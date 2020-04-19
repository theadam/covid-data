package main

import (
	"covid-tracker/data"
	"covid-tracker/fetch/jhu"
	"fmt"
	"sync"
	"time"
)

func runAction(wg *sync.WaitGroup, name string, action func()) {
    wg.Add(1)
	now := time.Now()

    go func () {
        action()
        fmt.Println("Finished " + name + " in " + time.Since(now).String())
        fmt.Println()
        wg.Done()
    }()
}

func main() {
	now := time.Now()
	fmt.Println("Loading all data")
	fmt.Println()

	points := jhu.GetData()

	var wg sync.WaitGroup

	runAction(&wg, "Writing World JSON data", func() { data.WriteWorldData(points) })
	runAction(&wg, "Writing Province JSON data", func() { data.WriteProvinceData(points) })
	runAction(&wg, "Writing State JSON data", func() { data.WriteStateData(points) })
	runAction(&wg, "Writing County JSON data", func() { data.WriteCountyData(points) })
	runAction(&wg, "Writing Date Range JSON data", func() { data.WriteDateRange(points) })
    wg.Wait()
    fmt.Println("Finished in " + time.Since(now).String())
}
