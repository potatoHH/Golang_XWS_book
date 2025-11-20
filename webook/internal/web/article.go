package web

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"

	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/service"
	ijwt "Book_Exp/webook/internal/web/jwt"
	"Book_Exp/webook/pkg/ginx"
	"Book_Exp/webook/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     service.ArticleService
	l       logger.LoggerV1
	intrSvc intrv1.InteractiveServiceClient
	biz     string
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1, intrSvc intrv1.InteractiveServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		biz:     "article",
		intrSvc: intrSvc,
	}
}

// 路由注册
func (a *ArticleHandler) RegisterRoutes(g *gin.Engine) {
	g.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	pub := g.Group("/pub")
	pub.POST("list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](a.List))
	pub.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))
	//点赞
	pub.POST("/like", ginx.WrapBodyAndToken[LikeReq, ijwt.UserClaims](a.Like))
	pub.POST("/cancel_like", ginx.WrapBodyAndToken[LikeReq, ijwt.UserClaims](a.Like))
	//收藏
	//pub.POST("/collect", ginx.WrapClaimsAndReq[CollectReq](a.Collect))

}

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq, uc ijwt.UserClaims) (ginx.Result, error) {
	_, err := a.intrSvc.Collect(ctx, &intrv1.CollectRequest{
		Biz:   a.biz,
		Uid:   uc.Uid,
		Cid:   req.Cid,
		BizId: req.Id,
	})
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc ijwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		_, err = a.intrSvc.Like(ctx, &intrv1.LikeRequest{
			Biz:   a.biz,
			Uid:   uc.Uid,
			BizId: req.Id,
		})
	} else {
		_, err = a.intrSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz:   a.biz,
			Uid:   uc.Uid,
			BizId: req.Id,
		})
	}
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

// 文章详情
func (a *ArticleHandler) Detail(ctx *gin.Context, usr ijwt.UserClaims) (ginx.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		//	ctx.JSON(200, ginx.Result{Code: 4, Msg: "参数错误"})
		//	a.l.Error("前端输入的ID不对", logger.Error(err))
		return ginx.Result{Code: 4, Msg: "参数错误"}, err
	}

	var eg errgroup.Group
	var art domain.Article
	eg.Go(func() error {
		uc := ctx.MustGet("users").(ijwt.UserClaims)
		art, err = a.svc.GetById(ctx, id, uc.Uid)
		return err
	})

	//这是不借助数据查询来判定的方法
	//if art.Author.Id != usr.Uid {
	//	//ctx.JSON(200, ginx.Result{Code: 4, Msg: "输入有误"})
	//	//如果公司有风控系统,这个时候就要上报这种非法访问的用户了
	//	a.l.Error("非法访问文章,创作者ID 不匹配", logger.Int64("uid", usr.Uid))
	//	return ginx.Result{Code: 4, Msg: "输入有误"}, nil

	//}
	//TODO  在这里获取文章的计数功能
	var getResp *intrv1.GetResponse
	eg.Go(func() error {
		uc := ctx.MustGet("users").(ijwt.UserClaims)
		getResp, err = a.intrSvc.Get(ctx, &intrv1.GetRequest{
			Biz:   a.biz,
			Uid:   uc.Uid,
			BizId: id,
		})
		if err != nil {
			//容忍错误的写法不返回 return
			a.l.Error("获取文章信息失败", logger.Error(err))
			//return ginx.Result{Code: 5, Msg: "系统错误"}, err
		}
		return err
	})
	//在这里等,要保证前面两个
	err = eg.Wait()
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	//TODO 增加阅读计数,这个功能 是不是可以让前端,主动刚发一个http请求 来增加一个计数
	go func() {
		//TODO 这里可以增加一个计数  go func  异步去执行
		_, er := a.intrSvc.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
			Biz:   a.biz,
			BizId: art.Id,
		})
		if er != nil {
			a.l.Error("增加阅读计数失败",
				logger.Error(err),
				logger.Int64("aid", art.Id),
			)
		}
	}()
	intr := getResp.Intr
	//这里不是借助数据库来判定的方法
	return ginx.Result{
		Data: ArtcleVO{
			Id:       art.Id,
			Title:    art.Title,
			Abstract: art.Abstrract(),
			//Content:  art.Content,
			//Author:   art.Author,
			Status:     art.Status.ToUnit8(),
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			Liked:      intr.Liked,
			Collected:  intr.Collected,
			CollentCnt: intr.CollectCnt,
			LikeCnt:    intr.LikeCnt,
			ReadCnt:    intr.ReadCnt,
		},
	}, nil

}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req struct {
		Id int64
	}
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
		a.l.Error("未发现用户的session 信息")
	}
	//检车输入,跳过
	//调用svc的代码
	err := a.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//打日志
		a.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "保存成功",
	})
}
func (a *ArticleHandler) Edit(ctx *gin.Context) {

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
		a.l.Error("未发现用户的session 信息")
	}
	//检车输入,跳过
	//调用svc的代码
	id, err := a.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//打日志
		a.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存成功",
		Data: id,
	})

}
func (a *ArticleHandler) Publish(ctx *gin.Context) {
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
		a.l.Error("未发现用户的session 信息")
	}
	//检车输入,跳过
	//调用svc的代码
	id, err := a.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//打日志
		a.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "保存成功",
		Data: id,
	})

}

func (a *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (ginx.Result, error) {
	res, err := a.svc.List(ctx, uc.Uid, req.Limit, req.Offset)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	//在列表页不显示 全文 只显示"摘要"
	return ginx.Result{
		Data: slice.Map[domain.Article, ArtcleVO](res,
			func(idx int, src domain.Article) ArtcleVO {
				return ArtcleVO{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstrract(),
					//Content:   src.Content,
					//Author: src.Author,
					Status: src.Status.ToUnit8(),
					Ctime:  src.Ctime.Format(time.DateTime),
					Utime:  src.Utime.Format(time.DateTime),
				}

			}),
	}, nil

}
