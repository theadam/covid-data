package data

import (
	"covid-tracker/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	. "github.com/ahmetb/go-linq"
)

var base = "./client/src/data/"

type Date struct {
	time.Time
}

func (a Date) CompareTo(bc Comparable) int {
	b := bc.(Date)
	if a.Before(b.Time) {
		return -1
	} else if b.Before(a.Time) {
		return 1
	}
	return 0
}

type DateData struct {
	Date      Date `json:"date"`
	Confirmed int  `json:"confirmed"`
	Deaths    int  `json:"deaths"`
}

const timeLayout = "2006-01-02"

func (u *Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.Format(timeLayout) + `"`), nil
}

func makeDateData(date Date, confirmed int, deaths int) DateData {
	return DateData{
		Date:      date,
		Confirmed: confirmed,
		Deaths:    deaths,
	}
}

func writeFile(json string, file string) {
	_, s, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(s))
	path := filepath.Join(root, base, file)
	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()
	_, err = f.WriteString(json)
	if err != nil {
		panic(err.Error())
	}
}

func DataToJson(data interface{}, points []*DataPoint) string {
	min, max := dateRange(points)
	value := reflect.ValueOf(data)

	for _, key := range value.MapKeys() {
		item := value.MapIndex(key)
		if !ensureRange(str(item.Interface()), utils.GetField("Dates", item.Elem().Interface()), min, max) {
			// Deletes item from map
			value.SetMapIndex(key, reflect.Value{})
		}
	}

	bytes, err := json.Marshal(value.Interface())
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}

type minMax struct {
	min time.Time
	max time.Time
}

func getDates(q Query) (*time.Time, *time.Time) {
	mm := q.AggregateWithSeed(
		minMax{min: time.Now().AddDate(1000, 0, 0)},
		func(acc interface{}, v interface{}) interface{} {
			times := acc.(minMax)
			date := utils.GetField("Date", v).(time.Time)

			if date.Before(times.min) {
				times.min = date
			}
			if date.After(times.max) {
				times.max = date
			}
			return times
		}).(minMax)
	return &mm.min, &mm.max
}

func dateRange(points []*DataPoint) (*time.Time, *time.Time) {
	return getDates(From(points))
}

func WriteDateRange(points []*DataPoint) {
	min, max := dateRange(points)
	cur := *min
	dates := make([]Date, 0)
	for !cur.After(*max) {
		dates = append(dates, Date{cur})
		cur = cur.AddDate(0, 0, 1)
	}
	data, err := json.Marshal(dates)
	if err != nil {
		panic(err.Error())
	}

	writeFile(string(data), "dateRange.json")
}

func str(item interface{}) string {
	bs, _ := json.Marshal(item)
	return string(bs)
}

func ensureRange(s string, items interface{}, min *time.Time, max *time.Time) bool {
	list := reflect.ValueOf(items)
	first := list.Index(0).Interface()
	minInList := utils.GetField("Date", first).(Date)
	if min.Before(minInList.Time) {
		if list.Len() == 1 {
			fmt.Println("Found what looks like a brand new case for " + s)
			return false
		}
		panic("Missing first date for " + s)
	}
	last := list.Index(list.Len() - 1).Interface()
	maxInList := utils.GetField("Date", last).(Date)
	if max.After(maxInList.Time) {
		panic("Missing last date for " + s)
	}
	return true
}

type Sums struct {
	Confirmed  int
	Deaths     int
	Population int
}

func getSums(query Query) Sums {
	return query.AggregateWithSeed(Sums{}, func(acc interface{}, v interface{}) interface{} {
		sums := acc.(Sums)
		value := reflect.Indirect(reflect.ValueOf(v))

		confirmed := int(value.FieldByName("Confirmed").Int())
		deaths := int(value.FieldByName("Deaths").Int())
		population := int(value.FieldByName("Population").Int())

		sums.Confirmed += confirmed
		sums.Deaths += deaths
		sums.Population += population
		return sums
	}).(Sums)
}

func grouping(groupShape interface{}) func(interface{}) interface{} {
	value := reflect.ValueOf(groupShape)
	typ := value.Type()
	return func(item interface{}) interface{} {
		result := reflect.New(typ).Elem()
		itemValue := reflect.Indirect(reflect.ValueOf(item))
		for i := 0; i < typ.NumField(); i++ {
			fname := typ.Field(i).Name
			field := itemValue.FieldByName(fname)
			result.FieldByName(fname).Set(field)
		}
		return result.Interface()
	}
}

func fromGroup(
	shape interface{},
) func(interface{}) interface{} {
	value := reflect.ValueOf(shape)
	typ := value.Type()
	return func(grp interface{}) interface{} {
		group := grp.(Group)
		groupKey := group.Key

		date := utils.GetField("Date", groupKey).(time.Time)
		sums := getSums(From(group.Group))

		sumsValue := reflect.ValueOf(sums)

		result := reflect.New(typ).Elem()

		for i := 0; i < typ.NumField(); i++ {
			fname := typ.Field(i).Name
			if fname == "Date" {
				result.FieldByName(fname).Set(reflect.ValueOf(Date{date}))
			} else if fname == "Confirmed" || fname == "Deaths" || fname == "Population" {
				result.FieldByName(fname).Set(sumsValue.FieldByName(fname))
			} else if !reflect.ValueOf(reflect.ValueOf(groupKey).FieldByName(fname)).IsZero() {
				result.FieldByName(fname).Set(reflect.ValueOf(utils.GetField(fname, groupKey)))
			}
		}
		return result.Interface()
	}
}

func aggregateTo(
	start interface{},
	getKey func(interface{}) interface{},
) (interface{}, func(interface{}, interface{}) interface{}) {
	value := reflect.ValueOf(start)
	typ := value.Type()
	elementType := typ.Elem().Elem()
	return start, func(acc interface{}, v interface{}) interface{} {
		objVal := reflect.ValueOf(acc)
		key := getKey(v)
		keyVal := reflect.ValueOf(key)

		ptrVal := objVal.MapIndex(keyVal)
		var val reflect.Value
		if reflect.ValueOf(ptrVal).IsZero() {
			ptrVal = reflect.New(elementType)
			val = ptrVal.Elem()
			for i := 0; i < elementType.NumField(); i++ {
				fname := elementType.Field(i).Name
				if fname == "Dates" {
					val.FieldByName(fname).Set(reflect.ValueOf(make([]DateData, 0)))
				} else {
					val.FieldByName(fname).Set(reflect.ValueOf(utils.GetField(fname, v)))
				}
			}
		} else {
			val = ptrVal.Elem()
		}

		newDates := append(utils.GetField("Dates", val.Interface()).([]DateData), makeDateData(
			utils.GetField("Date", v).(Date),
			utils.GetField("Confirmed", v).(int),
			utils.GetField("Deaths", v).(int),
		))

		ptrVal.Elem().FieldByName("Dates").Set(reflect.ValueOf(newDates))
		objVal.SetMapIndex(keyVal, ptrVal)
		return objVal.Interface()
	}
}

func CreateWorldData(points []*DataPoint) string {
	type shape struct {
		Date        Date   `json:"date"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		Confirmed   int    `json:"confirmed"`
		Deaths      int    `json:"deaths"`
		Population  int    `json:"population"`
	}
	type groupKey struct {
		Date        time.Time
		Country     string
		CountryCode string
	}
	type mapValue struct {
		Country     string     `json:"country"`
		CountryCode string     `json:"countryCode"`
		Dates       []DateData `json:"dates"`
		Population  int        `json:"population"`
	}

	obj := From(points).GroupBy(
		grouping(groupKey{}), utils.Id,
	).Select(fromGroup(
		shape{},
	)).OrderBy(
		utils.Field("Date"),
	).ThenBy(
		utils.Field("Country"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), utils.Field("CountryCode")))

	return DataToJson(obj, points)
}

func WriteWorldData(points []*DataPoint) {
	writeFile(CreateWorldData(points), "world.json")
}

func CreateProvinceData(points []*DataPoint) string {
	type shape struct {
		Date        Date   `json:"date"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		Province    string `json:"province"`
		Confirmed   int    `json:"confirmed"`
		Population  int    `json:"population"`
		Deaths      int    `json:"deaths"`
	}
	type groupKey struct {
		Date        time.Time
		Country     string
		CountryCode string
		Province    string
	}
	type mapValue struct {
		Country     string     `json:"country"`
		CountryCode string     `json:"countryCode"`
		Province    string     `json:"province"`
		Population  int        `json:"population"`
		Dates       []DateData `json:"dates"`
	}

	obj := From(points).Where(func(inter interface{}) bool {
        item := inter.(*DataPoint)
		return item.Province != "" && item.Country != "United States"
	}).GroupBy(
		grouping(groupKey{}), utils.Id,
	).Select(fromGroup(shape{})).OrderBy(
		utils.Field("Date"),
	).ThenBy(
		utils.Field("Country"),
	).ThenBy(
		utils.Field("Province"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), func(item interface{}) interface{} {
		return utils.GetField("CountryCode", item).(string) + "-" + utils.GetField("Province", item).(string)
	}))

	return DataToJson(obj, points)
}

func WriteProvinceData(points []*DataPoint) {
	writeFile(CreateProvinceData(points), "province.json")
}

func CreateStateData(points []*DataPoint) string {
	type shape struct {
		Province   string `json:"state"`
		FipsId     string `json:"fipsId"`
		Date       Date   `json:"date"`
		Population int    `json:"population"`
		Confirmed  int    `json:"confirmed"`
		Deaths     int    `json:"deaths"`
	}
	type groupKey struct {
		Date     time.Time
		Province string
	}
	type mapValue struct {
		Province   string     `json:"state"`
		FipsId     string     `json:"fipsId"`
		Population int        `json:"population"`
		Dates      []DateData `json:"dates"`
	}

	fipsMap := make(map[string]string)
	From(points).Where(func(item interface{}) bool {
		return utils.GetField("FipsId", item).(string) != ""
	}).DistinctBy(
		utils.Field("Province"),
	).Select(func(inter interface{}) interface{} {
		return KeyValue{
			Key:   utils.GetField("Province", inter),
			Value: utils.GetField("FipsId", inter).(string)[0:2],
		}
	}).ToMap(&fipsMap)

	obj := From(points).Where(func(item interface{}) bool {
        point := item.(*DataPoint)
		return point.Province != "" && point.Country == "United States"
	}).GroupBy(
		grouping(groupKey{}), utils.Id,
	).Select(
		fromGroup(shape{}),
	).Select(func(inter interface{}) interface{} {
		item := inter.(shape)
		fips, ok := fipsMap[item.Province]
		if ok {
			item.FipsId = fips
		}
		return item
    }).Where(func(inter interface{}) bool {
        point := inter.(shape)
        return point.FipsId != ""
    }).OrderBy(
		utils.Field("Date"),
	).ThenBy(
		utils.Field("Province"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), utils.Field("FipsId")))

	return DataToJson(obj, points)
}

func WriteStateData(points []*DataPoint) {
	writeFile(CreateStateData(points), "state.json")
}

func CreateCountyData(points []*DataPoint) string {
	type shape struct {
		Province      string `json:"state"`
		County     string `json:"county"`
		FipsId     string `json:"fipsId"`
		Date       Date   `json:"date"`
		Population int    `json:"population"`
		Confirmed  int    `json:"confirmed"`
		Deaths     int    `json:"deaths"`
	}
	type groupKey struct {
		Date   time.Time `json:"date"`
		Province  string    `json:"state"`
		County string    `json:"county"`
		FipsId string    `json:"fipsId"`
	}
	type mapValue struct {
		FipsId     string     `json:"id"`
		Population int        `json:"population"`
		Dates      []DateData `json:"dates"`
	}

	obj := From(points).Where(func(item interface{}) bool {
        point := item.(*DataPoint)
		return point.FipsId != "" && point.Country == "United States"
	}).GroupBy(
		grouping(groupKey{}), utils.Id,
	).Select(
		fromGroup(shape{}),
	).OrderBy(
		utils.Field("Date"),
	).ThenBy(
		utils.Field("Province"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), utils.Field("FipsId")))

	return DataToJson(obj, points)
}

func WriteCountyData(points []*DataPoint) {
	writeFile(CreateCountyData(points), "county.json")
}
