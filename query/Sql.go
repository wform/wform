package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/wform/wform/werror"
	"github.com/wform/wform/model"
	"github.com/wform/wform/utils"
)

// struct for build SQL
type Sql struct {
	column      string
	modelStruct model.ModelStruct
	join        string
	group       string
	having      string
	order       string
	limit       string
	offset      string
	option      string
	union       string
	updateSet   string

	rows []map[string]interface{}

	Cond Condition

	updateFieldMap map[string]bool
	omitFieldMap   map[string]bool
}

// update set in sql
func (sql *Sql) UpdateSet(updateSet string) *Sql {
	sql.updateSet = updateSet
	return sql
}

// fields in sql
func (sql *Sql) Select(column string) *Sql {
	sql.column = column
	return sql
}

// update field in sql
func (sql *Sql) Field(fields ...string) *Sql {
	updateFieldMap := map[string]bool{}
	for _, field := range fields {
		updateFieldMap[field] = true
	}
	sql.updateFieldMap = updateFieldMap
	return sql
}

// ignore update field in sql
func (sql *Sql) Omit(omit ...string) *Sql {
	omitFieldMap := map[string]bool{}
	for _, omt := range omit {
		omitFieldMap[omt] = true
	}
	sql.omitFieldMap = omitFieldMap
	return sql
}

/**
 * group by in sql
 */
func (sql *Sql) Group(group string) *Sql {
	sql.group = group
	return sql
}

// having condition in sql
func (sql *Sql) Having(having string) *Sql {
	sql.having = having
	return sql
}

// order by in sql
func (sql *Sql) Order(order string) *Sql {
	sql.order = order
	return sql
}

// limit in sql
func (sql *Sql) Limit(limit string) *Sql {
	sql.limit = limit
	return sql
}

// query offset in sql
func (sql *Sql) Offset(offset string) *Sql {
	sql.offset = offset
	return sql
}

// let join in sql
func (sql *Sql) Join(join string) *Sql {
	sql.join += " LEFT JOIN " + join
	return sql
}

// right join in sql
func (sql *Sql) RightJoin(join string) *Sql {
	sql.join += " RIGHT JOIN " + join
	return sql
}

// union in sql
func (sql *Sql) Union(sql2 *Sql) *Sql {
	sql.union += " UNION " + sql2.BuildSelect()
	return sql
}

// extra query (such as for update) in sql
func (sql *Sql) Option(option string) *Sql {
	sql.option = option
	return sql
}

// insert data rows or update data set in sql
func (sql *Sql) Rows(rows []map[string]interface{}) *Sql {
	sql.rows = rows
	return sql
}

// filter update columns
func (sql *Sql) FilterRows(rows []map[string]interface{}) *Sql {
	filterRows := []map[string]interface{}{}
	if len(sql.updateFieldMap) != 0 {
		if len(sql.omitFieldMap) == 0 {
			for _, row := range rows {
				filterRow := map[string]interface{}{}
				for field, val := range row {
					if _, exist := sql.updateFieldMap[field]; exist {
						filterRow[field] = val
					}
				}
				filterRows = append(filterRows, filterRow)
			}
		} else {
			for _, row := range rows {
				filterRow := map[string]interface{}{}
				for field, val := range row {
					_, fieldExist := sql.updateFieldMap[field]
					_, omitExist := sql.omitFieldMap[field]
					if fieldExist || !omitExist {
						filterRow[field] = val
					}
				}
				filterRows = append(filterRows, filterRow)
			}
		}

	} else if len(sql.omitFieldMap) != 0 {
		for _, row := range rows {
			filterRow := map[string]interface{}{}
			for field, val := range row {
				if _, exist := sql.omitFieldMap[field]; !exist {
					filterRow[field] = val
				}
			}
			filterRows = append(filterRows, filterRow)
		}
	} else {
		sql.rows = rows
	}
	return sql
}

// table in sql
func (sql *Sql) getTable() string {
	return sql.modelStruct.GetTableName()
}

func (sql *Sql) SetTable(tableName string) {
	sql.modelStruct.CustomizeTableName = tableName
}

// build select SQL statement
func (sql *Sql) BuildSelect() string {
	if sql.getTable() == "" {
		werror.WformPanic("table is empty")
	}
	queryString := "SELECT "
	if sql.column == "" {
		sql.column = "*"
	}
	queryString += sql.column + " FROM " + sql.getTable()

	if sql.join != "" {
		queryString += sql.join
	}

	if sql.HasCondition() {
		queryString += " WHERE " + sql.WhereSql()
	}
	if sql.group != "" {
		queryString += " GROUP BY " + sql.group
	}
	if sql.having != "" {
		queryString += " HAVING " + sql.having
	}
	if sql.order != "" {
		queryString += " ORDER BY " + sql.order
	}
	if sql.limit != "" {
		queryString += " LIMIT " + sql.limit
	}
	if sql.offset != "" {
		queryString += " OFFSET " + sql.offset
	}
	if sql.option != "" {
		queryString += " " + sql.option
	}
	if sql.union != "" {
		queryString += sql.union
	}
	return queryString
}

// Build SQL statement for delete record from table
func (sql *Sql) BuildDelete() string {
	if sql.getTable() == "" {
		werror.WformPanic("table is empty")
	}
	if !sql.HasCondition() {
		werror.WformPanic("condition is empty")
	}
	queryString := "DELETE FROM " + sql.getTable() + " WHERE " + sql.WhereSql()
	if sql.limit != "" {
		queryString += " LIMIT " + sql.limit
	}
	if sql.option != "" {
		queryString += " " + sql.option
	}
	return queryString
}

// build insert sql
func (sql *Sql) BuildInsert() string {
	queryString := "INSERT INTO " + sql.getTable()
	fields := []string{}
	rowValues := []string{}

	for field, _ := range sql.rows[0] {
		fields = append(fields, field)
	}

	for _, row := range sql.rows {
		rowStr := "("
		colsStr := []string{}
		for _, field := range fields {
			colsStr = append(colsStr, utils.SafeSqlValue(row[field]))
		}
		rowStr += strings.Join(colsStr, ",")
		rowStr += ")"
		rowValues = append(rowValues, rowStr)
	}
	if sql.option != "" {
		queryString += " " + sql.option
	}
	return queryString + "(" + strings.Join(fields, ",") + ") VALUES " + strings.Join(rowValues, ",")
}

// build update sql
func (sql *Sql) BuildUpdate() string {
	queryString := "UPDATE " + sql.getTable() + " SET "

	if sql.updateSet != "" {
		queryString += sql.updateSet
	} else {
		row := sql.rows[0]
		setSlice := []string{}
		for field, col := range row {
			switch col.(type) {
			case Exp:
				setSlice = append(setSlice, field+"="+col.(Exp).Value)
			default:
				setSlice = append(setSlice, field+"="+utils.SafeSqlValue(col))
			}
		}
		queryString += strings.Join(setSlice, ",")
	}
	queryString += " WHERE " + sql.WhereSql()
	if sql.option != "" {
		queryString += " " + sql.option
	}
	sql.updateSet = ""
	return queryString

}

// build where sql from complex condition
func (sql *Sql) WhereSql() string {
	var whereStr string
	for key, exp := range sql.Cond.exp {
		if key > 0 {
			if exp.logic == LogicOr {
				whereStr += " OR "
			} else if exp.logic == LogicAndNot {
				whereStr += " AND NOT ( "
			} else {
				whereStr += " AND "
			}
		} else if exp.logic == LogicAndNot {
			whereStr += " NOT ( "
		}

		switch exp.query.(type) {
		case string:
			whereStr += utils.ParsePrepare(exp.query.(string), exp.args...)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			whereStr += sql.modelStruct.PrimaryKey + "=" + fmt.Sprintf("%v", exp.query)
		case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []interface{}, []string:
			values := reflect.ValueOf(exp.query)
			var inStr []string
			for i := 0; i < values.Len(); i++ {
				inStr = append(inStr, fmt.Sprintf("%v", values.Index(i).Interface()))
			}
			whereStr += sql.modelStruct.PrimaryKey + " IN ( " + strings.Join(inStr, ",") + ")"
		case map[string]interface{}:
			var condStr  []string
			for field, value := range exp.query.(map[string]interface{}) {
				condStr = append(condStr, field+" = "+utils.SafeSqlValue(value))
			}
			whereStr += strings.Join(condStr, " AND ")
		case interface{}:
			attrMap := map[string]string{}
			modelStruct, _ := model.ParseStruct(exp.query)
			attrMap = modelStruct.FieldAttrMap
			objValue := reflect.ValueOf(exp.query)
			if objValue.Kind() == reflect.Struct {
				var condStr []string
				for field, attr := range attrMap {
					val := objValue.FieldByName(attr).Interface()
					if utils.ValIsEmpty(val) {
						continue
					}
					condStr = append(condStr, field+" = "+utils.SafeSqlValue(val))
				}
				whereStr += strings.Join(condStr, " AND ")
			} else {
				werror.WformPanic("type error")
			}
		default:
			werror.WformPanic(" where error ")
		}
		if exp.logic == LogicAndNot {
			whereStr += ")"
		}
	}
	return whereStr
}

// check if  the condition exist
func (sql *Sql) HasCondition() bool {
	return sql.Cond.hasCondition()
}

// set model
func (sql *Sql) SetModelStruct(modelStruct model.ModelStruct) {
	if sql.modelStruct.CustomizeTableName != "" {
		modelStruct.CustomizeTableName = sql.modelStruct.CustomizeTableName
	}
	sql.modelStruct = modelStruct
}
func (sql *Sql) GetModelStruct() model.ModelStruct {
	return sql.modelStruct
}

func (sql *Sql) GetRows() []map[string]interface{} {
	return sql.rows
}
