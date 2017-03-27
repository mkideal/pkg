package jsonx

import (
	"errors"
	"reflect"
)

// decodableValue
func decodableValue(ptr interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		t := reflect.TypeOf(ptr)
		if t == nil {
			return rv, errors.New("decode nil")
		}
		if t.Kind() != reflect.Ptr {
			return rv, errors.New("decode non-pointer " + t.String())
		}
		return rv, errors.New("decode nil " + t.String())
	}
	return rv, nil
}
