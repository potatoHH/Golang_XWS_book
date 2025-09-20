package ginx

import (
	"Book_Exp/webook/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//这个东西放到你们ginx插件库里面去, 技术含量不是很高,但是有心意

var L logger.LoggerV1

func WrapBody[T any](l logger.LoggerV1, fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		//TODO 下部分属于业务逻辑
		res, err := fn(ctx, req)
		if err != nil {
			//开始处理error,开始记录日志
			L.Error(
				"处理业务逻辑错误",
				logger.String("route", ctx.FullPath()),
				logger.String("path", ctx.Request.URL.Path),
				logger.Error(err),
			)
		}
		ctx.JSON(200, res)
	}
}

func WrapBodyAndToken[T any, C jwt.Claims](fn func(ctx *gin.Context, req T, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		//开始处理token
		var c C
		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		c, ok = val.(C) //断言
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		//TODO 下部分属于业务逻辑
		res, err := fn(ctx, req, c)
		if err != nil {
			//开始处理error,开始记录日志
			L.Error(
				"处理业务逻辑错误",
				logger.String("route", ctx.FullPath()),
				logger.String("path", ctx.Request.URL.Path),
				logger.Error(err),
			)
		}
		ctx.JSON(200, res)
	}

}

func WrapToken[C jwt.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//if err := ctx.Bind(&req); err != nil {
		//	return
		//}
	}
}
