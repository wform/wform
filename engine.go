package wform

import (
	"database/sql"

	"github.com/wform/wform/model"
	"github.com/wform/wform/query"
)

type Engine struct {
	dao interface{}
	sql *query.Sql

	db           *query.DB
	modelStruct  model.ModelStruct
	errorLogs    []error
	queryResults map[string]interface{}
}

// set the model for current query engine
func (engine *Engine) Dao(d interface{}) *Engine {
	if engine.sql == nil {
		engine.sql = &query.Sql{}
	}
	modelStruct, dao := model.ParseStruct(d)
	engine.dao = dao
	engine.sql.SetModelStruct(modelStruct)
	return engine
}

func (engine *Engine) judgeDao(d interface{}) {
	if engine.dao == nil {
		engine.Dao(d)
	}
}

func (engine *Engine) SetDb(d *sql.DB) {
	engine.db = &query.DB{
		SqlDb: d,
	}
}

// add error to db engine
func (engine *Engine) AddError(err error) {
	if err != nil {
		engine.errorLogs = append(engine.errorLogs, err)
	}
}

// Whether there is an error in query engine
func (engine *Engine) HasError() bool {
	return len(engine.errorLogs) > 0
}

// Get db error
func (engine *Engine) GetError(err error) []error {
	return engine.errorLogs
}

func Table(tableName string) *Engine {
	return E().Table(tableName)
}

func Where(query interface{}, values ...interface{}) *Engine {
	return E().Where(query, values...)
}

func Field(fields ...string) *Engine {
	return E().Field(fields...)
}

func Omit(omit ...string) *Engine {
	return E().Omit(omit...)
}

func Or(query interface{}, values ...interface{}) *Engine {
	return E().Or(query, values...)
}

func Not(query interface{}, values ...interface{}) *Engine {
	return E().Not(query, values...)
}

func Select(params ...interface{}) *Engine {
	return E().Select(params...)
}

func Group(group interface{}) *Engine {
	return E().Group(group)
}

func Having(having string) *Engine {
	return E().Having(having)
}

func Order(order interface{}) *Engine {
	return E().Order(order)
}

func Limit(limit interface{}) *Engine {
	return E().Limit(limit)
}

func Offset(offset interface{}) *Engine {
	return E().Offset(offset)
}

func Join(join string) *Engine {
	return E().Join(join)
}

func RightJoin(join string) *Engine {
	return E().RightJoin(join)
}

func Union(u *Engine) *Engine {
	return E().Union(u)
}

func Option(option string) *Engine {
	return E().Option(option)
}

func Find(dao interface{}) error {
	return E().Find(dao)
}

func Pluck(name string, fieldValues interface{}) error {
	return E().Pluck(name, fieldValues)
}

func Create(models ...interface{}) error {
	return E().Create(models...)
}

func Update(models ...interface{}) error {
	return E(models...).Update(models...)
}
func UpdateAll(models ...interface{}) error {
	return E(models...).UpdateAll(models...)
}
func Save(models ...interface{}) error {
	return E(models...).Save(models...)
}
func SaveAll(models ...interface{}) error {
	return E(models...).SaveAll(models...)
}

func Begin() *Engine {
	return E().Begin()
}

func InsertMany(rows []map[string]interface{}) error {
	return E().InsertMany(rows)
}

func Delete(params ...interface{}) error {
	return E().Delete(params...)
}
