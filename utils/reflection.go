package utils

import "reflect"

func Id(x interface{}) interface{} {
	return x
}

func GetField(name string, x interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(x)).FieldByName(name).Interface()
}

func Field(name string) func(interface{}) interface{} {
	return func(x interface{}) interface{} {
		return GetField(name, x)
	}
}
