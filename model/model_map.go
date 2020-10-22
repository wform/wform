package model

import (
	"reflect"

	"github.com/wform/wform/utils"
)

func Model2Map(obj interface{}, ignoreEmpty bool) map[string]interface{} {
	data := make(map[string]interface{})
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	data = parseEmbedStruct(data, objType, objValue, ignoreEmpty)
	return data
}

func parseEmbedStruct(data map[string]interface{}, objType reflect.Type, objValue reflect.Value, ignoreEmpty bool) map[string]interface{} {
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}
	if objType.Kind() == reflect.Struct {
		var fieldName string
		for i := 0; i < objType.NumField(); i++ {
			if objType.Field(i).Anonymous {
				data = parseEmbedStruct(data, objType.Field(i).Type, objValue.Field(i), ignoreEmpty)
			} else {
				fieldName = getWformTagFieldValue(objType.Field(i), "field")
				if fieldName == "" {
					fieldName = utils.ParseName(objType.Field(i).Name)
				}
				val := objValue.Field(i).Interface()
				if ignoreEmpty && utils.ValIsEmpty(val) {
					continue
				}
				data[fieldName] = val
			}
		}
	}
	return data
}
