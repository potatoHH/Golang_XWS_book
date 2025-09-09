package web

import "github.com/gin-gonic/gin"

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
}

// 路由注册
func (u *ArticleHandler) RegisterRoutes(service *gin.Engine) {

}
