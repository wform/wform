package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func ParseName(name string) string {
	bName := []byte(name)
	strLen := len(bName)
	var newName []byte
	for i := 0; i < strLen; i++ {
		curStrOrd := bName[i]
		if curStrOrd <= 91 {
			if i > 0 {
				newName = append(newName, 95)
			}
			newName = append(newName, bName[i]+32)
		} else {
			newName = append(newName, bName[i])
		}
	}
	return string(newName)
}

func ParsePrepare(query string, values ...interface{}) string {
	if len(values) == 0 {
		return query
	}
	var sprintValues  []interface{}
	for _, v := range values {
		switch v.(type) {
		case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
			values := reflect.ValueOf(v)
			inStr := []string{}
			for i := 0; i < values.Len(); i++ {
				inStr = append(inStr, fmt.Sprintf("%v", values.Index(i).Interface()))
			}
			query = strings.Replace(query, "?", "  (%v)", 1)
			sprintValues = append(sprintValues, strings.Join(inStr, ","))
		case []string, []interface{}:
			sprintValues = append(sprintValues, "'"+strings.Join(v.([]string), "','")+"'")
			query = strings.Replace(query, "?", " (%v)", 1)
		default:
			sprintValues = append(sprintValues, v)
			query = strings.Replace(query, "?", " '%v'", 1)
		}
	}
	return fmt.Sprintf(query, sprintValues...)
}

func ParseCommaFields(field interface{}) string {
	switch field.(type) {
	case string:
		return field.(string)
	case []string:
		return strings.Join(field.([]string), ",")
	}
	return ""
}
