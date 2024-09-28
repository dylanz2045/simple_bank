package db

import (
	"Project/utils"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestQueries *Queries
var TestDB *sql.DB
var err error

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("connot read config !")
		return
	}
	TestDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	TestQueries = New(TestDB)
	os.Exit(m.Run())
}
