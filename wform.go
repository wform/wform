package wform

import (
	"github.com/wform/wform/query"
	"github.com/wform/wform/utils"
	"github.com/wform/wform/werror"
)

func E(dao ...interface{}) *Engine {
	engine := Engine{}
	sqlDb, err := GetDb()
	if err != nil {
		werror.WformPanic(err)
	}
	engine.SetDb(sqlDb)
	if len(dao) != 0 {
		engine.Dao(dao[0])
	} else {
		engine.sql = &query.Sql{}
	}
	return &engine
}

/**
 * get query expression for update
 */
func Exp(expression string, params ...interface{}) query.Exp {
	return query.Exp{
		Value: utils.ParsePrepare(expression, params...),
	}
}
