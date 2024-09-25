package utils

import (
	"reflect"
	"strconv"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float32 | ~float64
}

func FormatFloat[V Float](value V) string {
	switch v := any(value).(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	default:
		panic("value type not float")
	}
}

func FormatInteger[V Integer](value V) string {
	switch v := any(value).(type) {
	case int8, int16, int32, int64, int:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	default:
		panic("value type not integer")
	}
}
