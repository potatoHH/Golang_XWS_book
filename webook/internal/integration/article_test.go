package integration

import (
	"Book_Exp/webook/internal/repository/dao/article"
	"Book_Exp/webook/internal/web"
	ijwt "Book_Exp/webook/internal/web/jwt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// 单元测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetUpSuite() {
	//在所有测试执行前,先执行一些内容
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("users", ijwt.UserClaims{
			Uid: 123,
		})
	})
	artHdl := web.NewArticleHandler()
	//注册好了路由
	artHdl.RegisterRoutes(s.server)

}

func (s *ArticleTestSuite) TestPublish() {}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		//集成测试
		before func(t *testing.T) // 集成测试前准备的数据
		after  func(t *testing.T) //集成测试后的验证数据
		//预期中的输入
		art      Article
		wantCode int           //http响应码
		wantRes  Result[int64] //希望http响应带上帖子的ID
	}{
		{
			name: "帖子保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "测试帖子",
					Content:  "测试帖子内容",
					AuthorId: 123,
				}, art)
			},
			art: Article{
				Title:   "测试帖子",
				Content: "测试帖子内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "OK",
			},
		},
		{
			name: "帖子修改成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					Ctime:    123,
					AuthorId: 123,
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "OK",
			},
		},
		{
			name: "修改其他人的帖子",
			before: func(t *testing.T) {
				//提前准备数据
				err := s.db.Create(article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 789,
					Utime:    234,
					Ctime:    123,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				var art article.Article
				err := s.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    123,
					Utime:    234,
					AuthorId: 789,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			//TODO 构造请求
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer(reqBody))
			require.NoError(t, err) //要求一定没有err
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")
			//这里就正常使用req
			resp := httptest.NewRecorder()
			//请求gin的路口  响应返回到resp
			s.server.ServeHTTP(resp, req)
			//断言结果
			assert.Equal(t, tc.wantCode, resp.Code)
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)

		})
	}
}

func (s *ArticleTestSuite) Testabc() {
	s.T().Log("这里是组合套件")
}
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `josn:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	//业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
