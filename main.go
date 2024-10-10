package main

import (
	"Project/api"
	db "Project/db/sqlc"
	"Project/utils"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {

	utils.Initlogger()
	logger := zap.L()
	if logger == nil {
		fmt.Println("logger is nil while it's shouldn't")
		os.Exit(-1)
	}
	config, err := utils.LoadConfig(".")
	if err != nil {
		logger.Error("cannot read config")
		return
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Error("cannot connect db")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		logger.Error("cannot listen Server")
	}
}
