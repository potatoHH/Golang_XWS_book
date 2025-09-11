package web

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/service"
	ijwt "Book_Exp/webook/internal/web/jwt"
	"Book_Exp/webook/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

// 路由注册
func (u *ArticleHandler) RegisterRoutes(service *gin.Engine) {
	service.Group("/articles")
	service.POST("/edit", u.Edit)
	service.POST("/publish", u.Publish)
}
func (u *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		//你可以考虑监控住这里
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		u.l.Error("未发现用户的session 信息")
	}
	//检车输入,跳过
	//调用svc的代码
	id, err := u.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//打日志
		u.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存成功",
		Data: id,
	})

}
func (u *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		//你可以考虑监控住这里
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		u.l.Error("未发现用户的session 信息")
	}
	//检车输入,跳过
	//调用svc的代码
	id, err := u.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//打日志
		u.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存成功",
		Data: id,
	})

}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `josn:"title"`
	Content string `json:"content"`
}

func (rep ArticleReq) toDomain(uid int64) domain.Article {
	var req ArticleReq
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}

}
