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
		logger.Sugar().Errorf("cannot read config for :%s", err)
		return
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Sugar().Errorf("cannot connect db for :%s", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		logger.Sugar().Errorf("failed : %s", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		logger.Sugar().Errorf("cannot listen Server for : %s", err)
	}
}
