package data

import (
	"database/sql"
	"fmt"
)

type DbConnector struct {
	db *sql.DB
}

func NewDbConnector(host, port, user, password, dbname string) *DbConnector {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Open():", err)
		panic(err)
	}

	return &DbConnector{db}
}

func (d *DbConnector) GetDb() *sql.DB {
	// TODO Change to pull of connections
	return d.db
}
