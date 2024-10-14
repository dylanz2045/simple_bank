package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpriesAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccess(ctx *gin.Context) {
	var req renewAccessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//需要验证这个token是否合法
	refreshPayload, err := server.token.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//验证这个会话是否被禁止
	if session.IsBlocked {
		err := errors.New("bolck session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	//验证这个会话，是否跟被验证的令牌是同一个请求
	if session.Username != refreshPayload.Username {
		err := errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err := errors.New("expried session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	//用于创建十分钟的身份验证token
	accessToken, accesspayload, err := server.token.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := renewAccessResponse{
		AccessToken:          accessToken,
		AccessTokenExpriesAt: accesspayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)

}
