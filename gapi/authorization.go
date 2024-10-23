package gapi

import (
	"Project/token"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBeaeer = "bearer"
)

// 用于验证GRPC的处理，让GRPC也知道现在是哪个用户在作函数的调用
func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing  authorization header")
	}

	//这里是获取到的标头：应该有两个部分组成
	//一个的这个令牌的类型
	//另一个是令牌的实体
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authorizationBeaeer {
		return nil, fmt.Errorf("unspported type")
	}

	accessToken := fields[1]
	accessPayload, err := server.token.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid accessToken")
	}
	return accessPayload, nil
}
