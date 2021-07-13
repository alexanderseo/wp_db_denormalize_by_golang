package main

import (
	"database/sql"
)
/*
Соединение с базой админки Wordpress
 */
func DbConnect() *sql.DB {
	db, err := sql.Open("mysql", "user:password@tcp(00.00.00.00:3306)/db_name")

	if err != nil {
		panic(err)
	}

	return db
}
