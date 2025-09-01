package jwt

import "github.com/gin-gonic/gin"

type Handler interface {
	setLoginToken(ctx *gin.Context, uid int64) error
	setRefreshToken(ctx *gin.Context, uid int64) error
	setLoginJwtToken(ctx *gin.Context, uid int64) error
	setRefreshJwtToken(ctx *gin.Context, uid int64) error
	setJwtToken(ctx *gin.Context, uid int64) error
	extractToken(ctx *gin.Context) string
}
