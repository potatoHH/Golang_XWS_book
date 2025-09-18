package pkg

import (
	"Book_Exp/webook/pkg/logger"

	"github.com/gin-gonic/gin"
)

//这个东西放到你们ginx插件库里面去, 技术含量不是很高,但是有心意

func WrapBody[T any](l logger.LoggerV1, fn func(req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		//TODO 下部分属于业务逻辑
		res, err := fn(req)
		if err != nil {
			//开始处理error,开始记录日志
			l.Error(
				"处理业务逻辑错误",
				logger.String("method", ctx.Request.Method),
				logger.String("path", ctx.Request.URL.Path),
				logger.Error(err),
			)
		}
		ctx.JSON
	}
}

type Result struct {
	//业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
