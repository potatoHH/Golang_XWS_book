package web

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// 确保Userhandler实现了handler的接口
// var _ handler = &UserHandler{}      //初始化了一个对象
var _ handler = (*UserHandler)(nil) //没有初始化任何对象  更优雅

const (
	userIdKey         = 1
	biz               = "login"
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// web里放的是跟路由相关的
type UserHandler struct {
	svc             service.UserServiceV1
	emilRegxExp     *regexp.Regexp
	passwordRegxExp *regexp.Regexp
	codeSvc         service.CodeServiceV1
}

func NewUserHandler(svc service.UserServiceV1, codeSvc service.CodeServiceV1) *UserHandler {
	return &UserHandler{
		emilRegxExp:     regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegxExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:             svc,
		codeSvc:         codeSvc,
	}
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) { // 注册路由
	//分组注册路由
	ug := server.Group("/users")
	ug.POST("/signup", c.Signup)
	ug.POST("/login", c.LoginJWT)
	ug.POST("/edit", c.Edit)
	ug.GET("/profile", c.Profile)
	ug.POST("/login_sms/code/send", c.SendLoginSmsCode)
}

// 路由接口

// 注册
func (c *UserHandler) Signup(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil { // 绑定参数
		return // 返回错误
	}
	isemail, err := c.emilRegxExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isemail {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	if req.Password != req.ConfirmPassword { // 判断密码是否一致
		ctx.String(http.StatusOK, "密码不一致")
		return
	}
	isPassword, err := c.passwordRegxExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码格式错误")
		return
	}
	// 实际调用服务层创建用户
	//调用一下 svc方法
	err = c.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicate { //邮箱重复
		ctx.String(http.StatusOK, "邮箱重复,请换一个邮箱")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误,注册失败")
		return
	}
	ctx.String(http.StatusOK, "注册成功")

} //

// 登录JWT
func (c *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	//参数绑定
	if err := ctx.Bind(&req); err != nil { // 绑定参数
		return // 返回错误
	}
	user, err := c.svc.Login(ctx, req.Email, req.Password) // 调用service
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return

	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if err = c.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")

	return

}

func (c *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		//设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //设置过期时间
		},
		Id:        uid,
		UserAgent: ctx.Request.UserAgent(), // 设置用户代理
	}
	//使用JWT设置登录状态  比如要求userid放入token中
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // 创建一个token 使用jwt
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		return err
	}

	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
func (c *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	//参数绑定
	if err := ctx.Bind(&req); err != nil { // 绑定参数
		return // 返回错误
	}
	user, err := c.svc.Login(ctx, req.Email, req.Password) // 调用service
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return

	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//登录成功之后,拿出session
	sess := sessions.Default(ctx) // 拿到session
	sess.Set("userId", user.Id)   // 设置session
	//sess.Set("update_time", user.Id)  //放在这里不太合适
	sess.Options(sessions.Options{ // 设置session的过期时间
		//Secure:   true,      // https  开发环境不要用
		//HttpOnly: true,      // js无法访问
		MaxAge: 30 * 60, // 表示30分钟
	})
	sess.Save() // 保存session

	ctx.String(http.StatusOK, "登录成功")
	return

}
func (c *UserHandler) logOut(ctx *gin.Context) {
	sess := sessions.Default(ctx) // 拿到session
	sess.Options(sessions.Options{ // 设置session的过期时间
		//Secure:   true,      // https  开发环境不要用
		//HttpOnly: true,      // js无法访问
		MaxAge: -1, // 表示立即删除或清除这个 session cookie

	})
	sess.Save() // 保存session
	ctx.String(http.StatusOK, "登录成功")
	return

}
func (c *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		//TODO 注意 其他字段 尤其是密码,邮箱和手机号,修改都要通过别的字段 邮箱和手机都要验证 &&密码
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		//在这里尝试校验
		if req.Nickname == "" {
			ctx.JSON(http.StatusOK, Result{
				Code: 4,
				Msg:  "昵称不能为空",
			})
			return
		}
		if len(req.AboutMe) > 1024 {
			ctx.JSON(http.StatusOK, Result{
				Code: 4,
				Msg:  "关于我太长",
			})
			return
		}
		birthday, err := time.Parse(time.DateOnly, req.Birthday)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 4,
				Msg:  "时间格式错误",
			})
			return
		}
		uc := ctx.MustGet("user").(UserClaims)
		err = c.svc.UpdateNonSensitiveInfo(ctx, domain.User{
			Id:       uc.Id,
			Nickname: req.Nickname,
			Birthday: birthday,
			AboutMe:  req.AboutMe,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		ctx.JSON(http.StatusOK, Result{
			Msg: "Ok",
		})
	}
}
func (c *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	u, err := c.svc.Profile(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: u.Email,
	})

}

// ProfileJWT 用户详情 jwt版本
func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	uc := ctx.MustGet("user").(UserClaims)
	u, err := c.svc.Profile(ctx, uc.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Format(time.DateOnly),
		AboutMe:  u.AboutMe,
	})
}
func (c *UserHandler) SendLoginSmsCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
	}
	err = c.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码发送次数太多,请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

}
func (c *UserHandler) LoginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := c.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码错误",
		})
		return
	}
	//我这个手机号,会不会是一个新用户呢?
	//
	user, err := c.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//这里怎么办,从哪里来
	if err = c.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证码校验成功",
	})

}

type UserClaims struct { // jwt
	jwt.RegisteredClaims // 继承
	//声明你自己要放进去token里面的数据
	Id int64

	UserAgent string
}
