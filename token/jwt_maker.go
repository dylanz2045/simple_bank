package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size:must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// 创建一个新的令牌
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	//创建一个新的Token负载
	payload, err := NewPayload(username, duration)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//创建一个token令牌
	jwtTocker := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtTocker.SignedString([]byte(maker.secretKey))
}

// 验证这个令牌是否合法
// 检查令牌的签名方法、使用密钥验证签名、以及解析令牌中的声明
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	//定义的一个匿名函数：作为提供验证JWT令牌需要的秘钥
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
