package data

import (
	"encoding/json"
	"testing"
	"time"
)

var now = time.Now()

var sampleWorld = []DataPoint{
    DataPoint{
        Province: "Quebec",
        Country: "Canada",
        CountryCode: "124",
        Confirmed: 10,
        Deaths: 9,
        Population: 400,
        Date: now.AddDate(0, 0, -1),
    },
    DataPoint{
        Province: "Quebec",
        Country: "Canada",
        CountryCode: "124",
        Confirmed: 10,
        Deaths: 9,
        Population: 400,
        Date: now,
    },
}

var sampleUs = []CountyData{
    CountyData{
        State: "Georgia",
        County: "Fulton",
        FipsId: "12345",
        Confirmed: 10,
        Deaths: 9,
        Population: 400,
        Date: now.AddDate(0, 0, -1),
    },
    CountyData{
        State: "Georgia",
        County: "Fulton",
        FipsId: "12345",
        Confirmed: 10,
        Deaths: 9,
        Population: 400,
        Date: now,
    },
}

func parseJson(js string) map[string]interface{} {
    var res map[string]interface{}
    json.Unmarshal([]byte(js), &res)
    return res
}

func TestMissingDatePanics(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("The code was expected to panic")
        }
    }()
    CreateWorldData(append(sampleWorld, DataPoint{
        Province: "",
        Country: "China",
        CountryCode: "123",
        Confirmed: 10,
        Deaths: 9,
        Population: 400,
        Date: time.Now(),
    }), []CountyData{})
}


func TestWorldCombinePoints(t *testing.T) {
    data := parseJson(CreateWorldData(sampleWorld, []CountyData{}))
    canada := data["124"].(map[string]interface{})
    if len(data) != 1 {
        t.Errorf("Too many items in map: %d", len(data))
    }
    if (canada["countryCode"] != "124") {
        t.Errorf("Unexpected country code: %s", canada["countryCode"])
    }
    if (canada["country"] != "Canada") {
        t.Errorf("Unexpected country: %s", canada["country"])
    }
    if (canada["population"] != 400.0) {
        t.Errorf("Unexpected population: %s", canada["population"])
    }
}

func TestDates(t *testing.T) {
    data := parseJson(CreateWorldData(sampleWorld, []CountyData{}))
    canada := data["124"].(map[string]interface{})
    dates := canada["dates"].([]interface{})
    if (len(dates) != 2) {
        t.Errorf("Canada has the wrong number of dates: %d", len(dates))
    }
}

func TestState(t *testing.T) {
    data := parseJson(CreateStateData([]DataPoint{}, sampleUs))
    _, ok := data["12"]
    if !ok {
        t.Error("Georgia's fips id is not found")
    }
    if sampleUs[0].FipsId != "12345" {
        t.Error("FipsId was changed")
    }
}

func test(t *testing.T) {
    data := parseJson(CreateStateData([]DataPoint{}, sampleUs))
    _, ok := data["12"]
    if !ok {
        t.Error("Georgia's fips id is not found")
    }
}
