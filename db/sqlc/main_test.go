package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "user=postgres password=cst4Ever dbname=postgres host=localhost port=5432 sslmode=disable"
)

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	TestQueries = New(conn)
	os.Exit(m.Run())

}
