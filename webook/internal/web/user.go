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
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// web里放的是跟路由相关的
type UserHandler struct {
	svc             *service.UserService
	emilRegxExp     *regexp.Regexp
	passwordRegxExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emilRegxExp:     regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegxExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:             svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) { // 注册路由
	//分组注册路由
	ug := server.Group("/users")
	ug.POST("/login", u.LoginJWT)
	ug.POST("/signup", u.Signup)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)

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
	if err == service.ErrUserDuplicateEmail { //邮箱重复
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
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	//参数绑定
	if err := ctx.Bind(&req); err != nil { // 绑定参数
		return // 返回错误
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password) // 调用service
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不对")
		return

	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	claims := UserClaims{
		Uid: user.Id,
	}
	//使用JWT设置登录状态  比如要求userid放入token中
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // 创建一个token 使用jwt
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)

	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return

}
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	//参数绑定
	if err := ctx.Bind(&req); err != nil { // 绑定参数
		return // 返回错误
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password) // 调用service
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

func (u *UserHandler) logOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)  // 拿到session
	sess.Options(sessions.Options{ // 设置session的过期时间
		//Secure:   true,      // https  开发环境不要用
		//HttpOnly: true,      // js无法访问
		MaxAge: -1, // 表示立即删除或清除这个 session cookie

	})
	sess.Save() // 保存session
	ctx.String(http.StatusOK, "登录成功")
	return

}
func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的profile")

}

func (u *UserHandler) ProfileJwt(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
}

type UserClaims struct { // jwt
	jwt.RegisteredClaims
	//声明你自己要放进去token里面的数据
	Uid int64
}
