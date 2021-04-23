package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

func rowExist(query string, args ...interface{}) (bool, error) {
	var exist bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("rowExists: %v", err)
	}
	return exist, nil
}

// initializes the database. panics if a failure.
func init() {
	var err error
	// do NOT use in production
	connStr := "user=postgres password=? dbname=ecocosts sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}
