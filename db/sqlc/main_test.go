package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testStore Store

const (
	dbDriver = "postgres"
	dbSource = "user=userdb password=cst4Ever dbname=ailanzbase host=8.134.97.76 port=5432 "
)

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(conn)
	os.Exit(m.Run())

}
