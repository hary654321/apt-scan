package utils

import (
	"reflect"
	"strconv"
	"strings"
)

func ToMap(param interface{}) map[string]string {
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	m := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		// 通过interface方法来获取key所对应的值
		var cell string
		switch s := v.Field(i).Interface().(type) {
		case string:
			cell = s
		case []string:
			cell = strings.Join(s, "; ")
		case int:
			cell = strconv.Itoa(s)
		case Stringer:
			cell = s.String()
		default:
			continue
		}
		m[t.Field(i).Name] = cell
	}
	return m
}

type Stringer interface {
	String() string
}
