package orm

import (
	"reflect"
	"sync"
)

type TableReflector func(interface{}) Table

var (
	reflectedTables   map[reflect.Type]TableReflector
	reflectedTablesMu sync.RWMutex
)

func typeof(v interface{}) reflect.Type {
	typ := reflect.TypeOf(v)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func reflectTable(v interface{}) (Table, error) {
	if table, ok := v.(Table); ok {
		return table, nil
	}
	if rotable, ok := v.(ReadonlyTable); ok {
		return extendReadonlyTable(rotable)
	}
	if wotable, ok := v.(WriteonlyTable); ok {
		return extendWriteonlyTable(wotable)
	}
	if tableinfo, ok := v.(TableInfo); ok {
		return extendTableInfo(tableinfo)
	}
	return nil, nil
}

func reflectReadonlyTable(v interface{}) (ReadonlyTable, error) {
	return nil, nil
}

func extendReadonlyTable(rotable ReadonlyTable) (Table, error) {
	return nil, nil
}

func extendWriteonlyTable(wotable WriteonlyTable) (Table, error) {
	return nil, nil
}

func extendTableInfo(tableinfo TableInfo) (Table, error) {
	return nil, nil
}
