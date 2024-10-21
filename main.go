package main

import (
	"Project/api"
	db "Project/db/sqlc"
	"Project/gapi"
	"Project/pb"
	"Project/utils"
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	_ "Project/doc/statik"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	runDBMigration(config.MigrationURL, config.DBSource)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runDBMigration(migrationURL string, dbSource string) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("logger is nil while it's shouldn't")
		os.Exit(-1)
	}
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Sugar().Errorf("can not create new migrate instance : %s --->main.go 61", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Sugar().Errorf("can not run migration up : %s --->main.go 64", err)
	}
	logger.Info("db migrated successful")

}
func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		logger.Sugar().Errorf("failed : %s --->main.go 55", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		logger.Sugar().Errorf("cannot listen Server for : %s --->main.go 60", err)
	}
}

func runGrpcServer(config utils.Config, store db.Store) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("logger is nil while it's shouldn't")
		os.Exit(-1)
	}
	server, err := gapi.NewServer(config, store)
	if err != nil {
		logger.Sugar().Errorf("failed : %s --->main.go 72", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		logger.Error("can not create Listener----> runGrpcServer")
	}
	logger.Sugar().Info("start GRPC server at %s----> run GRPCServer", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Error("can not start Server----> runGrpcServer")
	}

}

func runGatewayServer(config utils.Config, store db.Store) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("logger is nil while it's shouldn't")
		os.Exit(-1)
	}
	server, err := gapi.NewServer(config, store)
	if err != nil {
		logger.Sugar().Errorf("failed : %s --->main.go 99", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		logger.Sugar().Errorf("can not register handle server for :%s --->main.go 115", err)
		os.Exit(-1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//创建一个监听处理器，用于监听端口时，响应的一个网页静态文件
	statikFS, err := fs.New()
	if err != nil {
		logger.Sugar().Errorf("can not creatik fs for : %s --->main.go 124", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	//创建http的监听端口
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		logger.Sugar().Errorf("can not create listener for :%s", err)
	}
	logger.Sugar().Info("start HTTP gateway server at %s----> run GatewayServer", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		logger.Sugar().Errorf("can not create server for :%s --->main.go 136", err)
	}

}
