package main

import (
	"Project/api"
	db "Project/db/sqlc"
	"Project/gapi"
	"Project/pb"
	"Project/utils"
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logger *zap.Logger

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

	runGrpcServer(config, store)

}
func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		logger.Sugar().Errorf("failed : %s", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		logger.Sugar().Errorf("cannot listen Server for : %s", err)
	}
}

func runGrpcServer(config utils.Config, store db.Store) {
	utils.Initlogger()
	logger := zap.L()
	if logger == nil {
		fmt.Println("logger is nil while it's shouldn't")
		os.Exit(-1)
	}
	server, err := gapi.NewServer(config, store)
	if err != nil {
		logger.Sugar().Errorf("failed : %s", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		logger.Error("can not create Listener----> runGrpcServer")
	}
	logger.Info(" Listening----> run FRPCServer")
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Error("can not start Server----> runGrpcServer")
	}

}
