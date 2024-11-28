package db

import (
	//  "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DBManager struct {
	DB *sqlx.DB
}

func NewDBConnection(databaseURL string) *DBManager {
	dbx, err := sqlx.Open("sqlite3", "popupbase.db")

	if err != nil {
		panic(err)
	}

	driver, err := sqlite3.WithInstance(dbx.DB, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"ql", driver)
	m.Up()

	return &DBManager{
		DB: dbx,
	}
}
