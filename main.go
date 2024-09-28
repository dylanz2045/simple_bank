package main

import (
	"Project/api"
	db "Project/db/sqlc"
	"Project/utils"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("connot load config", err)
		return
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connnot start server:", err)
	}
}
