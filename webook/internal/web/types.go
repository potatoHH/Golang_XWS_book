package web

import (
	"Book_Exp/webook/pkg/ginx"

	"github.com/gin-gonic/gin"
)

type handler interface {
	RegisterRoutes(server *gin.Engine)
}
type Result1 = ginx.Result
