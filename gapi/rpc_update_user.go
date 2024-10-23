package gapi

import (
	db "Project/db/sqlc"
	"Project/pb"
	"Project/utils"
	"Project/val"
	"context"
	"database/sql"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPay, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	//验证传输过来的信息是否符合规范
	if violations := validateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentEror(violations)
	}
	//验证信息携带的令牌是否为自己的
	if authPay.Username != req.Username {
		return nil, permissionDenyError()
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "that user hasn't register :%s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot find that user:%s", err)
	}
	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
	}
	if req.Password != nil {
		err := utils.CheckPassword(*req.Password, user.HashedPassword)
		if err == nil {
			return nil, status.Error(codes.PermissionDenied, "the password is same as before!")
		}
		hashedPassword, err := utils.HashedPassword(*req.Password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password:%s ", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "that user hasn't register :%s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot upodate that user:%s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: ConverUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.FullName != nil {
		if err := val.ValidateFullname(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("fullname", err))
		}
	}
	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	return violations
}
