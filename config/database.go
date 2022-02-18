package config

import (
	"database/sql"
	"os"
)

func DbConn2(dbname string) (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "rooty"
	dbPass := "P@ssw0rd"
	dbName := "ejlog3"

	if string(dbname) != "" {
		dbName = dbname
	}

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(172.26.193.96:3306)/"+dbName+"?timeout=2s")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(100)

	return db
}

func DbConn(dbname string) (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := string(os.Getenv("DB_USER"))
	dbPass := string(os.Getenv("DB_PASSWORD"))
	dbHost := string(os.Getenv("DB_HOST"))
	dbPort := string(os.Getenv("DB_PORT"))
	dbName := "ejlog3"

	if string(dbname) != "" {
		dbName = dbname
	}

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?timeout=2s")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
