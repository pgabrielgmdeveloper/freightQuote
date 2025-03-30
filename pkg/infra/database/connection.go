package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewDbInstance(dbDriver, dbHost, dbPort, dbName, sslMode, dbPassword, dbUser string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
