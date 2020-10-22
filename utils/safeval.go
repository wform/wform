package utils

import (
	"fmt"
	"strings"
)

// add slashes to query string for special character
func SafeSqlValue(val interface{}) string {
	var safeVal string
	switch val.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, int, uint:
		safeVal = fmt.Sprintf("%v", val)
	default:
		safeVal = "'" + strings.Replace(fmt.Sprintf("%v", val), "'", "\\'", -1) + "'"
	}
	return safeVal
}

func ValIsEmpty(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float64, float32:
		if fmt.Sprint(val) == "0" {
			return true
		}
	case string:
		if val == "" {
			return true
		}
	default:
		return true

	}
	return false
}
