package gapi

import (
	db "Project/db/sqlc"
	"Project/pb"
	"Project/utils"
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "that user hasn't register :%s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot find that user:%s", err)
	}
	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password is wrong :%s", err)
	}
	//用于创建十分钟的身份验证token
	accessToken, accesspayload, err := server.token.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not create Token :%s", err)
	}
	//下面用于创建一天的Session
	refreshToken, refreshpayload, err := server.token.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not create Token :%s", err)
	}
	session, err := server.store.CreateSession(ctx,
		db.CreateSessionParams{
			ID:           refreshpayload.ID,
			Username:     refreshpayload.Username,
			RefreshToken: refreshToken,
			UserAgent:    "",
			ClientIp:     "",
			IsBlocked:    false,
			ExpiresAt:    refreshpayload.ExpiredAt,
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not create Session :%s", err)
	}
	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accesspayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshpayload.ExpiredAt),
		User:                  ConverUser(user),
	}
	return rsp, nil
}
