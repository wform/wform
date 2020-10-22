package query

import (
	"database/sql"
	"fmt"
)

type DB struct {
	SqlDb *sql.DB
	SqlTx *sql.Tx
	HasTx bool
}

func (db *DB) Query(Sql string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(Sql)
	if db.HasTx {
		return db.SqlTx.Query(Sql, args...)
	} else {
		return db.SqlDb.Query(Sql, args...)
	}
}

func (db *DB) Exec(Sql string, args ...interface{}) (sql.Result, error) {
	if db.HasTx {
		return db.SqlTx.Exec(Sql, args...)
	} else {
		return db.SqlDb.Exec(Sql, args...)
	}
}

func (db *DB) Begin() error {
	tx, err := db.SqlDb.Begin()
	db.HasTx = true
	db.SqlTx = tx
	return err
}

func (db *DB) Commit() error {
	err := db.SqlTx.Commit()
	db.HasTx = false
	return err
}

func (db *DB) Rollback() error {
	err := db.SqlTx.Rollback()
	db.HasTx = false
	return err
}
