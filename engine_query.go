package worm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/wform/worm/werror"

	"github.com/wform/worm/model"
	"github.com/wform/worm/query"
	"github.com/wform/worm/utils"
)

// Change query table
func (engine *Engine) Table(tableName string) *Engine {
	engine.sql.SetTable(tableName)
	return engine
}

func (engine *Engine) Where(query interface{}, values ...interface{}) *Engine {
	engine.sql.Cond.Where(query, values...)
	return engine
}

func (engine *Engine) Field(fields ...string) *Engine {
	engine.sql.Field(fields...)
	return engine
}

func (engine *Engine) Omit(omit ...string) *Engine {
	engine.sql.Omit(omit...)
	return engine
}

func (engine *Engine) Or(query interface{}, values ...interface{}) *Engine {
	engine.sql.Cond.Or(query, values...)
	return engine
}

func (engine *Engine) Not(query interface{}, values ...interface{}) *Engine {
	engine.sql.Cond.Not(query, values...)
	return engine
}

func (engine *Engine) Select(params ...interface{}) *Engine {
	if len(params) == 0 {
		werror.WormPanic("Select params error")
	}
	switch params[0].(type) {
	case string:
		engine.sql.Select(utils.ParsePrepare(params[0].(string), params[1:]...))
	case []string:
		engine.sql.Select(strings.Join(params[0].([]string), ","))
	default:
		werror.WormPanic("Select params error")
	}
	return engine
}

func (engine *Engine) Group(group interface{}) *Engine {
	engine.sql.Group(utils.ParseCommaFields(group))
	return engine
}

func (engine *Engine) Having(having string) *Engine {
	engine.sql.Having(having)
	return engine
}

func (engine *Engine) Order(order interface{}) *Engine {
	engine.sql.Order(utils.ParseCommaFields(order))
	return engine
}

func (engine *Engine) Limit(limit interface{}) *Engine {
	engine.sql.Limit(fmt.Sprintf("%v", limit))
	return engine
}

func (engine *Engine) Offset(offset interface{}) *Engine {
	engine.sql.Offset(fmt.Sprintf("%v", offset))
	return engine
}

func (engine *Engine) Join(join string) *Engine {
	engine.sql.Join(join)
	return engine
}

func (engine *Engine) RightJoin(join string) *Engine {
	engine.sql.RightJoin(join)
	return engine
}

func (engine *Engine) Union(u *Engine) *Engine {
	engine.sql.Union(u.sql)
	return engine
}

func (engine *Engine) Option(option string) *Engine {
	engine.sql.Option(option)
	return engine
}

func (engine *Engine) Clear() {
	modelStruct := engine.sql.GetModelStruct()
	engine.sql = &query.Sql{}
	engine.sql.SetModelStruct(modelStruct)
}

func (engine *Engine) One() (interface{}, error) {
	engine.sql.Limit("1")
	result, err := engine.All()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

func (engine *Engine) All() ([]interface{}, error) {
	SQL := engine.sql.BuildSelect()
	defer engine.Clear()
	rows, err := engine.db.Query(SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dao := engine.dao
	fieldMap := engine.sql.GetModelStruct().FieldAttrMap

	daoValue := reflect.ValueOf(dao)
	var daoIsPtr bool
	if daoValue.Kind() == reflect.Ptr {
		daoValue = daoValue.Elem()
		daoIsPtr = true
	}

	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))

	for index, _ := range columns {
		var ignored interface{}
		values[index] = &ignored
	}

	result := []interface{}{}
	for rows.Next() {
		newDaoPtr := reflect.New(daoValue.Type())
		newDao := reflect.ValueOf(newDaoPtr.Interface())
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for i, column := range columns {
			if attr, exist := fieldMap[column]; exist {
				pval := values[i].(*interface{})
				val := *pval
				if val == nil {
					continue
				}

				var varStr string
				switch val.(type) {
				case time.Time:
					varStr = val.(time.Time).String()
				case []byte:
					varStr = string(val.([]byte)[:])
				}

				var modelAttr reflect.Value
				if newDao.Kind() == reflect.Ptr {
					modelAttr = newDao.Elem().FieldByName(attr)
				} else {
					modelAttr = newDao.FieldByName(attr)
				}
				switch k := modelAttr.Kind(); k {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					valI, err := strconv.Atoi(varStr)
					if err != nil {
						continue
					}
					modelAttr.SetInt(int64(valI))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					valI, err := strconv.Atoi(varStr)
					if err != nil {
						continue
					}
					modelAttr.SetUint(uint64(valI))
				case reflect.Float64, reflect.Float32:
					valF, err := strconv.ParseFloat(varStr, 64)
					if err != nil {
						continue
					}
					modelAttr.SetFloat(valF)
				case reflect.String:
					modelAttr.SetString(varStr)
				}
			}
		}
		if daoIsPtr {
			result = append(result, newDao.Interface())
		} else {
			result = append(result, newDao.Elem().Interface())
		}
	}
	return result, nil
}

// TODO optimize for
func (engine *Engine) Find(dao interface{}) error {
	engine.Dao(dao)
	daoValue := reflect.ValueOf(dao)
	if daoValue.Kind() != reflect.Ptr {
		werror.WormPanic("address is nil")
	}
	if daoValue.IsNil() {
		werror.WormPanic("value is nil")
	}
	//addressDaoValue := daoValue
	daoValue = daoValue.Elem()
	if daoValue.Kind() == reflect.Slice {
		results, err := engine.All()
		if err != nil {
			return err
		}

		for _, v := range results {
			daoValue.Set(reflect.Append(daoValue, reflect.ValueOf(v)))
		}
	} else {
		foundDao, err := engine.One()
		if err != nil {
			return err
		}
		if foundDao != nil {
			daoValue.Set(reflect.ValueOf(foundDao).Elem())
		}
	}
	return nil
}

// @TODO optimize for
func (engine *Engine) Pluck(name string, fieldValues interface{}) error {
	rfVal := reflect.ValueOf(fieldValues)
	if rfVal.Kind() != reflect.Ptr {
		werror.WormPanic("fieldValues must be pointer")
	}
	rfVal = rfVal.Elem()

	switch rfVal.Kind() {
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Slice:
	case reflect.String:
	default:
		werror.WormPanic("FieldValues must be slice or array")
	}

	if rfVal.Kind() != reflect.Slice {
		engine = engine.Limit(1)
	}

	engine.Select(name)
	SQL := engine.sql.BuildSelect()
	engine.Select("")
	rows, err := engine.db.Query(SQL)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		switch rfVal.Kind() {
		case reflect.Slice:
			paramPtr := reflect.New(rfVal.Type().Elem()).Interface()
			paramValue := reflect.ValueOf(paramPtr).Elem()
			switch k := rfVal.Type().Elem().Kind(); k {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var valI int64
				err := rows.Scan(&valI)
				if err != nil {
					return err
				}
				paramValue.SetInt(valI)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var valI uint64
				err := rows.Scan(&valI)
				if err != nil {
					return err
				}
				paramValue.SetUint(valI)
			case reflect.Float64, reflect.Float32:
				var valF float64
				err := rows.Scan(&valF)
				if err != nil {
					return err
				}
				paramValue.SetFloat(valF)
			case reflect.String:
				var fieldStr string
				err := rows.Scan(&fieldStr)
				if err != nil {
					return err
				}
				paramValue.SetString(fieldStr)
			default:
				werror.WormPanic(fmt.Sprintf("does not support this type %v", k))
			}
			rfVal.Set(reflect.Append(rfVal, paramValue))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var valI int64
			err := rows.Scan(&valI)
			if err != nil {
				return err
			}
			rfVal.SetInt(int64(valI))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var valI uint64
			err := rows.Scan(&valI)
			if err != nil {
				return err
			}
			rfVal.SetUint(uint64(valI))
		case reflect.Float64, reflect.Float32:
			var valF float64
			err := rows.Scan(&valF)
			if err != nil {
				return err
			}
			rfVal.SetFloat(valF)
		case reflect.String:
			var fieldStr string
			err := rows.Scan(&fieldStr)
			if err != nil {
				return err
			}
			rfVal.SetString(fieldStr)
		}
	}

	return nil
}

// @TODO test
func (engine *Engine) Rows() (rows *sql.Rows, err error) {
	SQL := engine.sql.BuildSelect()
	rows, err = engine.db.Query(SQL)
	engine.AddError(err)
	return
}

// Create data list from models
func (engine *Engine) Create(models ...interface{}) error {
	if len(models) == 0 {
		if engine.dao != nil {
			models = []interface{}{engine.dao}
		} else {
			werror.WormPanic(" model can not be nil")
		}
	} else {
		engine.Dao(models[0])
	}
	for _, mod := range models {
		modVal := reflect.ValueOf(mod)
		if modVal.Kind() != reflect.Ptr {
			werror.WormPanic("created model params must be pointer")
		}
		method := modVal.MethodByName("BeforeCreate")
		if method.Kind() == reflect.Func {
			method.Call([]reflect.Value{modVal})
		}

		modelMap := model.Model2Map(mod, false)
		modelStruct := engine.sql.GetModelStruct()
		_, exist := modelMap[modelStruct.PrimaryKey]
		if exist {
			delete(modelMap, modelStruct.PrimaryKey)
		}
		engine.sql.Rows([]map[string]interface{}{modelMap})
		result, err := engine.query(engine.insertSql())
		if err != nil {
			engine.AddError(err)
			return err
		}
		if engine.sql.GetModelStruct().PrimaryKey != "" {
			lastInsertId, err := result.LastInsertId()
			if err != nil {
				werror.WormPanic(err)
			}
			fieldName := modVal.Elem().FieldByName(engine.sql.GetModelStruct().GetPrimaryKeyAttr())
			switch fieldName.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldName.SetUint(uint64(lastInsertId))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldName.SetInt(lastInsertId)
			case reflect.String:
				fieldName.SetString(strconv.Itoa(int(lastInsertId)))
			default:
				werror.WormPanic("Primary Key " + engine.sql.GetModelStruct().GetPrimaryKeyAttr() + ", the kind " + fieldName.Kind().String() + " not support Primary Key")
			}

		}
		method = modVal.MethodByName("AfterCreate")
		if method.Kind() == reflect.Func {
			method.Call([]reflect.Value{modVal})
		}
	}
	return nil
}

func (engine *Engine) Update(params ...interface{}) error {
	return engine.updateToDb(true, params...)
}

func (engine *Engine) UpdateAll(params ...interface{}) error {
	return engine.updateToDb(false, params...)
}

// update query
func (engine *Engine) updateToDb(ignoreEmpty bool, params ...interface{}) error {
	if len(params) == 0 {
		if engine.dao != nil {
			params = []interface{}{engine.dao}
		} else {
			werror.WormPanic(" model can not be nil")
		}
	}

SAVE_DONE:
	for key, param := range params {
		modVal := reflect.ValueOf(param)
		if modVal.Kind() == reflect.Ptr {
			modVal = modVal.Elem()
		}
		var SQL string
		switch modVal.Kind() {
		case reflect.String:
			if key == 0 {
				setStr := param.(string)
				
				if len(params) >= 2 {
					setStr = ""
					updateMap := map[string]interface{}{}
					for i := 1; i<len(params); i += 2 {
						lastParam := params[i-1]
						switch lastParam.(type) {
						case string:
							updateMap[lastParam.(string)] = params[i]
						}

					}
					engine.sql.FilterRows([]map[string]interface{}{updateMap})
				}
				engine.sql.UpdateSet(setStr)
				SQL = engine.updateSql()
				result, err := engine.query(SQL)
				engine.setRowsAffected(result)
				if err != nil {
					return err
				}
				break SAVE_DONE
			}
		case reflect.Struct:
			engine.judgeDao(params[0])
			engine.sql.FilterRows([]map[string]interface{}{model.Model2Map(param, ignoreEmpty)})
			SQL = engine.updateSql()
		case reflect.Map:
			modMap, assertResult := param.(map[string]interface{})
			if assertResult {
				break
			}
			engine.sql.FilterRows([]map[string]interface{}{modMap})
			SQL = engine.updateSql()
		}
		if SQL == "" {
			werror.WormPanic("params type error")
		}
		result, err := engine.query(SQL)
		engine.setRowsAffected(result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (engine *Engine) Save(models ...interface{}) error {
	return engine.saveToDb(true, models...)
}

func (engine *Engine) SaveAll(models ...interface{}) error {
	return engine.saveToDb(false, models...)
}

func (engine *Engine) saveToDb(ignoreEmpty bool, models ...interface{}) error {
	if len(models) == 0 {
		if engine.dao != nil {
			models = []interface{}{engine.dao}
		} else {
			werror.WormPanic(" model can not be nil")
		}
	}

	for _, mod := range models {
		modVal := reflect.ValueOf(mod)
		var modValAddr reflect.Value
		if modVal.Kind() == reflect.Ptr {
			modValAddr = modVal
			modVal = modVal.Elem()
		}
		var SQL string

		if modVal.Kind() == reflect.Struct {
			engine.Dao(mod)
		}
		modelStruct := engine.sql.GetModelStruct()

		var sqlObj *query.Sql
		var modelMap map[string]interface{}
		var isInsert bool

		if modelStruct.PrimaryKey == "" {
			werror.WormPanic(fmt.Sprintf("%v has no PrimaryKey tag", mod))
		}

		switch modVal.Kind() {
		case reflect.Struct:
			modelMap = model.Model2Map(mod, ignoreEmpty)
			sqlObj = engine.sql.Rows([]map[string]interface{}{modelMap})

			primaryVal, exist := modelMap[modelStruct.PrimaryKey]
			var hasPrimaryVal bool
			if exist && !utils.ValIsEmpty(primaryVal) {
				hasPrimaryVal = true
			}
			if hasPrimaryVal {
				oldCond := sqlObj.Cond.GetExp()
				sqlObj.Cond.Where(modelMap[modelStruct.PrimaryKey])
				engine.sql = sqlObj
				SQL = engine.updateSql()
				sqlObj.Cond.Empty()
				sqlObj.Cond.SetExp(oldCond)
				engine.sql = sqlObj
			} else {
				if exist {
					delete(modelMap, modelStruct.PrimaryKey)
				}

				engine.sql = sqlObj
				SQL = engine.insertSql()
				isInsert = true
			}

		case reflect.Map:
			var assertResult bool
			modelMap, assertResult = mod.(map[string]interface{})
			if !assertResult {
				werror.WormPanic(fmt.Sprintf("%v type error", mod))
			}
			sqlObj = engine.sql.Rows([]map[string]interface{}{modelMap})
			if _, exist := modelMap[modelStruct.PrimaryKey]; exist {
				sqlObj.Cond.Where(modelMap[modelStruct.PrimaryKey])
				engine.sql = sqlObj
				SQL = engine.updateSql()
			} else {
				SQL = engine.insertSql()
				isInsert = true
			}
		}

		method := modVal.MethodByName("BeforeSave")
		if method.Kind() == reflect.Func {
			method.Call([]reflect.Value{modVal})
		}
		result, err := engine.query(SQL)
		if err != nil {
			return err
		}

		method = modVal.MethodByName("AfterSave")
		if method.Kind() == reflect.Func {
			method.Call([]reflect.Value{modVal})
		}

		if isInsert && modValAddr.Kind() == reflect.Ptr && engine.sql.GetModelStruct().PrimaryKey != "" {
			lastInsertId, err := result.LastInsertId()
			if err != nil {
				werror.WormPanic(err)
			}
			fieldName := modValAddr.Elem().FieldByName(engine.sql.GetModelStruct().GetPrimaryKeyAttr())
			switch fieldName.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldName.SetUint(uint64(lastInsertId))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldName.SetInt(lastInsertId)
			case reflect.String:
				fieldName.SetString(strconv.Itoa(int(lastInsertId)))
			default:
				werror.WormPanic("Primary Key " + engine.sql.GetModelStruct().GetPrimaryKeyAttr() + ", the kind " + fieldName.Kind().String() + " not support Primary Key")
			}
		}

		engine.setRowsAffected(result)
	}
	return nil
}

func (engine *Engine) Begin() *Engine {
	err := engine.db.Begin()
	if err != nil {
		engine.AddError(err)
	}
	return engine
}

func (engine *Engine) Commit() *Engine {
	err := engine.db.Commit()
	if err != nil {
		engine.AddError(err)
		engine.Rollback()
	}
	return engine
}

func (engine *Engine) Rollback() *Engine {
	err := engine.db.Rollback()
	engine.AddError(err)
	return engine
}

// batch insert data list
func (engine *Engine) InsertMany(rows []map[string]interface{}) error {
	engine.sql.Rows(rows)
	SQL := engine.insertSql()
	result, err := engine.query(SQL)
	if err != nil {
		engine.AddError(err)
		return err
	}
	engine.setRowsAffected(result)
	return nil
}

// delete from where clause
func (engine *Engine) Delete(params ...interface{}) error {

	useDelete := false
	if len(params) >= 2 {
		var assertResult bool
		useDelete, assertResult = params[1].(bool)
		if !assertResult {
			werror.WormPanic(fmt.Sprintf("%v must be bool", params[1]))
		}
	}

	if len(params) >= 1 {
		modelStruct, _ := model.ParseStruct(params[0])
		mod := reflect.ValueOf(params[0])
		engine.Dao(params[0]).Where(fmt.Sprintf("%v=%v", modelStruct.PrimaryKey, mod.FieldByName(modelStruct.PrimaryKey).Interface()))
	}

	var deleteSql string
	if useDelete {
		deleteSql = engine.deleteSql()
	} else {
		if engine.sql.GetModelStruct().DeletedAtColumn == "" {
			deleteSql = engine.deleteSql()
		} else {
			err := engine.Update(map[string]interface{}{
				engine.sql.GetModelStruct().DeletedAtColumn: time.Now().Unix(),
			})
			if err != nil {
				return err
			}
		}
	}
	result, err := engine.query(deleteSql)
	engine.Clear()
	if err != nil {
		engine.AddError(err)
		return err
	}
	engine.setRowsAffected(result)
	return nil
}

func (engine *Engine) setRowsAffected(result driver.Result) {
	rowsAffected, err := result.RowsAffected()
	engine.AddError(err)
	if err == nil {
		if engine.queryResults == nil {
			engine.queryResults = map[string]interface{}{}
		}
		engine.queryResults["RowsAffected"] = rowsAffected
	}
}

func (engine *Engine) RowsAffected() int64 {
	var rowsAffected int64
	if row, exist := engine.queryResults["RowsAffected"]; exist {
		rowsAffected = row.(int64)
	}
	return rowsAffected
}

func (engine *Engine) query(sql string) (sql.Result, error) {
	fmt.Println(sql)
	return engine.db.Exec(sql)
}

func (engine *Engine) insertSql() string {
	if queryCall["insertCall"] != nil {
		return queryCall["insertCall"](engine.sql)
	}
	return engine.sql.BuildInsert()
}

func (engine *Engine) updateSql() string {
	if queryCall["updateCall"] != nil {
		return queryCall["updateCall"](engine.sql)
	}
	return engine.sql.BuildUpdate()
}

func (engine *Engine) deleteSql() string {
	if queryCall["deleteCall"] != nil {
		return queryCall["deleteCall"](engine.sql)
	}
	return engine.sql.BuildDelete()
}

var queryCall = make(map[string]func(*query.Sql) string)

func Register(name string, qc func(*query.Sql) string) {
	queryCall[name] = qc
}
