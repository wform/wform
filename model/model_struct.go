package model

import (
	"reflect"
	"strings"
	"sync"

	"github.com/wform/worm/utils"
	"github.com/wform/worm/werror"
)

// the struct attributes
type ModelStruct struct {
	TableName          string
	CustomizeTableName string
	PrimaryKey         string
	DeletedAtColumn    string
	FieldAttrMap       map[string]string
	SubTagMap          map[string]map[string]string
}

func (m ModelStruct) GetPrimaryKeyAttr() string {
	return m.FieldAttrMap[m.PrimaryKey]
}

func (m ModelStruct) GetTableName() string {
	if m.CustomizeTableName != "" {
		return m.CustomizeTableName
	} else {
		return m.TableName
	}
}

var modelStructParseMap sync.Map

func ParseStruct(d interface{}) (ModelStruct, interface{}) {
	structType := reflect.TypeOf(d)
	var daoIsPtr bool
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() == reflect.Slice {
		structType = structType.Elem()
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
			daoIsPtr = true
		}
	} else {
		daoIsPtr = true
	}
	if structType.Kind() != reflect.Struct {
		werror.WormPanic("struct error")
	}
	structValue := reflect.New(structType)
	if !daoIsPtr {
		structValue = structValue.Elem()
	}
	daoValue := structValue

	packagePath := structType.PkgPath()
	structName := structType.Name()
	var mapKey string
	if packagePath != "" && structName != "" {
		mapKey = packagePath + "/" + structName
		if daoIsPtr {
			mapKey = "*" + mapKey
		}
	}
	var parsedModelStruct ModelStruct
	var exist bool

	// Load from cache
	if mapKey != "" {
		var modStruct interface{}
		modStruct, exist = modelStructParseMap.Load(mapKey)
		if exist {
			parsedModelStruct = modStruct.(ModelStruct)
			parsedModelVal, _ := modelStructParseMap.Load(mapKey + "_val")
			return parsedModelStruct, parsedModelVal
		}
	}
	fieldAttrMap := map[string]string{}
	subTagMap := map[string]map[string]string{}
	_, methodExist := structType.MethodByName("TableName")
	tableName := ""

	if methodExist {
		rtnValue := structValue.MethodByName("TableName").Call(nil)
		if rtnValue != nil {
			tableName = rtnValue[0].String()
		}
	}
	if tableName == "" {
		tableName = utils.ParseName(structType.Name())
	}

	var columnName, primaryKey, deletedAtColumn string
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Type.Kind() == reflect.Struct {
			innerType := structType.Field(i).Type
			for j := 0; j < innerType.NumField(); j++ {
				fieldTagMap := getWormTagMap(innerType.Field(j))
				columnName = fieldTagMap["column"]
				if columnName == "" {
					columnName = utils.ParseName(innerType.Field(j).Name)
				}
				fieldAttrMap[columnName] = innerType.Field(j).Name
				subTagMap[columnName] = fieldTagMap
				_, exist := fieldTagMap["primary_key"]
				if exist {
					primaryKey = columnName
				}

				_, exist = fieldTagMap["deleted_at"]
				if exist {
					deletedAtColumn = columnName
				}
			}
			continue
		}
		fieldTagMap := getWormTagMap(structType.Field(i))
		columnName = fieldTagMap["column"]
		if columnName == "" {
			columnName = utils.ParseName(structType.Field(i).Name)
		}
		firstChar := structType.Field(i).Name[:1]
		if strings.ToUpper(firstChar) == firstChar {
			fieldAttrMap[columnName] = structType.Field(i).Name
		}
		subTagMap[columnName] = fieldTagMap
		_, exist := fieldTagMap["primary_key"]
		if exist {
			primaryKey = columnName
		}

		_, exist = fieldTagMap["deleted_at"]
		if exist {
			deletedAtColumn = columnName
		}
	}
	parsedModelStruct.TableName = tableName
	parsedModelStruct.PrimaryKey = primaryKey
	parsedModelStruct.DeletedAtColumn = deletedAtColumn
	parsedModelStruct.FieldAttrMap = fieldAttrMap
	parsedModelStruct.SubTagMap = subTagMap
	if mapKey != "" {
		modelStructParseMap.Store(mapKey, parsedModelStruct)
		modelStructParseMap.Store(mapKey+"_val", daoValue.Interface())
	}
	return parsedModelStruct, daoValue.Interface()
}

func getWormTagMap(params ...interface{}) map[string]string {
	if len(params) == 0 {
		werror.WormPanic("params error")
	}
	field := params[0].(reflect.StructField)
	var tagKey string
	if len(params) > 1 {
		tagKey = params[1].(string)
	}
	wormTag := field.Tag.Get("db")
	fieldTagMap := map[string]string{}
	if wormTag != "" {
		wormTagList := strings.Split(wormTag, ";")
		for _, subTag := range wormTagList {
			subTagSplit := strings.Split(subTag, ":")
			switch len(subTagSplit) {
			case 2:
				fieldTagMap[subTagSplit[0]] = subTagSplit[1]
			default:
				fieldTagMap[subTagSplit[0]] = ""
			}
			if tagKey != "" {
				return fieldTagMap
			}
		}
	}
	return fieldTagMap
}

func getWormTagFieldValue(field reflect.StructField, fieldName string) string {
	fieldTagMap := getWormTagMap(field, fieldName)
	return fieldTagMap[fieldName]
}
