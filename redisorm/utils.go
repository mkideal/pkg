package orm

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrFieldNotFound    = errors.New("field not found")
	ErrUnexpectedLength = errors.New("unexpected length")
)

func ToString(v interface{}) string {
	switch value := v.(type) {
	case int:
		return strconv.FormatInt(int64(value), 10)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(int64(value), 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(uint64(value), 10)
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return fmt.Sprintf("%v", v)
	}
}
