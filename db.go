package worm

import (
	"database/sql"
)

var dbInstance *sql.DB

var lastSQLConfig []string

// OpenDb open the database connection
func OpenDb(driverName string, dbDN string) (*sql.DB, error) {
	lastSQLConfig = []string{driverName, dbDN}
	var err error
	dbInstance, err = sql.Open(lastSQLConfig[0], lastSQLConfig[1])
	return dbInstance, err
}

// GetDb get the database connection object
func GetDb() (*sql.DB, error) {
	var err error
	if dbInstance == nil {
		dbInstance, err = OpenDb(lastSQLConfig[0], lastSQLConfig[1])
	}
	return dbInstance, err
}
