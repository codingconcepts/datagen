package runner

import (
	"database/sql"
	"log"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

var (
	db   *sql.DB
	mock sqlmock.Sqlmock
)

func resetMock() {
	var err error
	if db, mock, err = sqlmock.New(); err != nil {
		log.Fatalf("error creating sqlmock: %v", err)
	}
}
