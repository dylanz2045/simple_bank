package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PateoMaker struct {
	paseto        *paseto.V2
	symmertricKey []byte
}

func NewPasetoMaker(symmertricKey string) (Maker, error) {
	if len(symmertricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size : must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PateoMaker{
		paseto:        paseto.NewV2(),
		symmertricKey: []byte(symmertricKey),
	}
	return maker, nil
}

// 创建一个新的Paseto的令牌
// 一个作为类的方法，那个调用的实体是可以取的
func (pateoMaker *PateoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token,err := pateoMaker.paseto.Encrypt(pateoMaker.symmertricKey, payload, nil)
	return token , payload , err
}

func (pateoMaker *PateoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := pateoMaker.paseto.Decrypt(token, pateoMaker.symmertricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
