package gapi

import (
	db "Project/db/sqlc"
	"Project/pb"
	"Project/token"
	"Project/utils"
	"fmt"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config utils.Config
	token  token.Maker
	store  db.Store
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}
	server := &Server{
		config: config,
		store:  store,
		token:  tokenMaker,
	}
	return server, nil
}
