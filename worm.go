package worm

import (
	"github.com/wform/worm/query"
	"github.com/wform/worm/utils"
	"github.com/wform/worm/werror"
)

func E(dao ...interface{}) *Engine {
	engine := Engine{}
	sqlDb, err := GetDb()
	if err != nil {
		werror.WormPanic(err)
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
