package data

import (
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

func DataToJson(data interface{}, world []DataPoint, us []CountyData) string {
	min, max := dateRange(world, us)
	value := reflect.ValueOf(data)

	for _, key := range value.MapKeys() {
		item := value.MapIndex(key)
		if !ensureRange(str(item.Interface()), getField("Dates", item.Elem().Interface()), min, max) {
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

func getDates(q Query) (time.Time, time.Time) {
	mm := q.AggregateWithSeed(
		minMax{min: time.Now().AddDate(1000, 0, 0)},
		func(acc interface{}, v interface{}) interface{} {
			times := acc.(minMax)
			date := getField("Date", v).(time.Time)

			if date.Before(times.min) {
				times.min = date
			}
			if date.After(times.max) {
				times.max = date
			}
			return times
		}).(minMax)
	return mm.min, mm.max
}

func dateRange(points []DataPoint, counties []CountyData) (time.Time, time.Time) {
	min, max := getDates(From(points))
	cmin, cmax := getDates(From(counties))

	if min.After(cmin) {
		min = cmin
	}
	if max.Before(cmax) {
		max = cmax
	}
	return min, max
}

func WriteDateRange(points []DataPoint, counties []CountyData) {
	min, max := dateRange(points, counties)
	cur := min
	dates := make([]Date, 0)
	for !cur.After(max) {
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

func ensureRange(s string, items interface{}, min time.Time, max time.Time) bool {
	list := reflect.ValueOf(items)
	first := list.Index(0).Interface()
	minInList := getField("Date", first).(Date)
	if min.Before(minInList.Time) {
		if list.Len() == 1 {
			fmt.Println("Found what looks like a brand new case for " + s)
			return false
		}
		panic("Missing first date for " + s)
	}
	last := list.Index(list.Len() - 1).Interface()
	maxInList := getField("Date", last).(Date)
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
		value := reflect.ValueOf(v)

		confirmed := int(value.FieldByName("Confirmed").Int())
		deaths := int(value.FieldByName("Deaths").Int())
		population := int(value.FieldByName("Population").Int())

		sums.Confirmed += confirmed
		sums.Deaths += deaths
		sums.Population += population
		return sums
	}).(Sums)
}

func id(x interface{}) interface{} {
	return x
}

func getField(name string, x interface{}) interface{} {
	return reflect.ValueOf(x).FieldByName(name).Interface()
}

func field(name string) func(interface{}) interface{} {
	return func(x interface{}) interface{} {
		return reflect.ValueOf(x).FieldByName(name).Interface()
	}
}

func grouping(groupShape interface{}) func(interface{}) interface{} {
	value := reflect.ValueOf(groupShape)
	typ := value.Type()
	return func(item interface{}) interface{} {
		result := reflect.New(typ).Elem()
		itemValue := reflect.ValueOf(item)
		for i := 0; i < typ.NumField(); i++ {
			fname := typ.Field(i).Name
			field := itemValue.FieldByName(fname)
			result.FieldByName(fname).Set(field)
		}
		return result.Interface()
	}
}

func fromGroupWithOverrides(
	shape interface{},
	getKey func(interface{}) interface{},
	overrides map[interface{}]map[time.Time]Sums,
) func(interface{}) interface{} {
	value := reflect.ValueOf(shape)
	typ := value.Type()
	return func(grp interface{}) interface{} {
		group := grp.(Group)
		groupKey := group.Key
		overrideKey := getKey(groupKey)

		date := getField("Date", groupKey).(time.Time)
		override, ok := overrides[overrideKey]
		var sums Sums
		if ok {
			sums = override[date]
		} else {
			sums = getSums(From(group.Group))
		}

		sumsValue := reflect.ValueOf(sums)

		result := reflect.New(typ).Elem()

		for i := 0; i < typ.NumField(); i++ {
			fname := typ.Field(i).Name
			if fname == "Date" {
				result.FieldByName(fname).Set(reflect.ValueOf(Date{date}))
			} else if fname == "Confirmed" || fname == "Deaths" || fname == "Population" {
				result.FieldByName(fname).Set(sumsValue.FieldByName(fname))
			} else if !reflect.ValueOf(reflect.ValueOf(groupKey).FieldByName(fname)).IsZero() {
				result.FieldByName(fname).Set(reflect.ValueOf(getField(fname, groupKey)))
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
					val.FieldByName(fname).Set(reflect.ValueOf(getField(fname, v)))
				}
			}
		} else {
			val = ptrVal.Elem()
		}

		newDates := append(getField("Dates", val.Interface()).([]DateData), makeDateData(
			getField("Date", v).(Date),
			getField("Confirmed", v).(int),
			getField("Deaths", v).(int),
		))

		ptrVal.Elem().FieldByName("Dates").Set(reflect.ValueOf(newDates))
		objVal.SetMapIndex(keyVal, ptrVal)
		return objVal.Interface()
	}
}

func fromGroup(shape interface{}) func(interface{}) interface{} {
	return fromGroupWithOverrides(shape, id, make(map[interface{}]map[time.Time]Sums))
}

func CreateWorldData(world []DataPoint, us []CountyData) string {
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

	usMap := make(map[time.Time]Sums)
	From(us).GroupBy(
		field("Date"),
		id,
	).Select(func(group interface{}) interface{} {
		return KeyValue{
			Key:   group.(Group).Key,
			Value: getSums(From(group.(Group).Group)),
		}
	}).ToMap(&usMap)

	obj := From(world).GroupBy(
		grouping(groupKey{}), id,
	).Select(fromGroupWithOverrides(
		shape{},
		field("Country"),
		map[interface{}]map[time.Time]Sums{
			"United States": usMap,
		},
	)).OrderBy(
		field("Date"),
	).ThenBy(
		field("Country"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), field("CountryCode")))

	return DataToJson(obj, world, us)
}

func WriteWorldData(world []DataPoint, us []CountyData) {
	writeFile(CreateWorldData(world, us), "world.json")
}

func CreateProvinceData(world []DataPoint, us []CountyData) string {
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

	obj := From(world).Where(func(item interface{}) bool {
		return getField("Province", item).(string) != ""
	}).GroupBy(
		grouping(groupKey{}), id,
	).Select(fromGroup(shape{})).OrderBy(
		field("Date"),
	).ThenBy(
		field("Country"),
	).ThenBy(
		field("Province"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), func(item interface{}) interface{} {
		return getField("CountryCode", item).(string) + "-" + getField("Province", item).(string)
	}))

	return DataToJson(obj, world, us)
}

func WriteProvinceData(world []DataPoint, us []CountyData) {
	writeFile(CreateProvinceData(world, us), "province.json")
}

func CreateStateData(world []DataPoint, us []CountyData) string {
	type shape struct {
		State      string `json:"state"`
		FipsId     string `json:"fipsId"`
		Date       Date   `json:"date"`
		Population int    `json:"population"`
		Confirmed  int    `json:"confirmed"`
		Deaths     int    `json:"deaths"`
	}
	type groupKey struct {
		Date  time.Time
		State string
	}
	type mapValue struct {
		State      string     `json:"state"`
		FipsId     string     `json:"fipsId"`
		Population int        `json:"population"`
		Dates      []DateData `json:"dates"`
	}

	type fipsGroup struct {
		State string
	}

	fipsMap := make(map[string]string)
	From(us).Where(func(item interface{}) bool {
		return getField("FipsId", item).(string) != ""
	}).DistinctBy(
        field("State"),
    ).Select(func(inter interface{}) interface {} {
        return KeyValue{
            Key: getField("State", inter),
            Value: getField("FipsId", inter).(string)[0:2],
        }
    }).ToMap(&fipsMap)

    obj := From(us).Where(func(item interface{}) bool {
		return getField("State", item).(string) != ""
	}).GroupBy(
		grouping(groupKey{}), id,
	).Select(
		fromGroup(shape{}),
    ).Select(func (inter interface{}) interface{} {
        item := inter.(shape)
        fips, ok := fipsMap[item.State]
        if ok {
            item.FipsId = fips
        }
        return item
    }).OrderBy(
		field("Date"),
	).ThenBy(
		field("State"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), field("FipsId")))

	return DataToJson(obj, world, us)
}

func WriteStateData(world []DataPoint, us []CountyData) {
	writeFile(CreateStateData(world, us), "state.json")
}

func CreateCountyData(world []DataPoint, us []CountyData) string {
	type shape struct {
		State      string `json:"state"`
		County     string `json:"county"`
		FipsId     string `json:"fipsId"`
		Date       Date   `json:"date"`
		Population int    `json:"population"`
		Confirmed  int    `json:"confirmed"`
		Deaths     int    `json:"deaths"`
	}
	type groupKey struct {
		Date   time.Time   `json:"date"`
		State  string `json:"state"`
		County string `json:"county"`
		FipsId string `json:"fipsId"`
	}
	type mapValue struct {
		FipsId         string     `json:"id"`
		Population int        `json:"population"`
		Dates      []DateData `json:"dates"`
	}

	obj := From(us).Where(func(item interface{}) bool {
		return getField("FipsId", item).(string) != ""
	}).GroupBy(
		grouping(groupKey{}), id,
	).Select(
		fromGroup(shape{}),
	).OrderBy(
		field("Date"),
	).ThenBy(
		field("State"),
	).AggregateWithSeed(aggregateTo(make(map[string]*mapValue), field("FipsId")))

	return DataToJson(obj, world, us)
}

func WriteCountyData(world []DataPoint, us []CountyData) {
	writeFile(CreateCountyData(world, us), "county.json")
}
