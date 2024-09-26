package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestQueries *Queries
var TestDB *sql.DB
var err error

const (
	dbDriver = "postgres"
	dbSource = "user=postgres password=cst4Ever dbname=mypostgres host=localhost port=5432 sslmode=disable"
)

func TestMain(m *testing.M) {
	fmt.Println("setup test")
	TestDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	TestQueries = New(TestDB)
	os.Exit(m.Run())
}
