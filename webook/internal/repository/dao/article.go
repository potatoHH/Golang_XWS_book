package dao

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}
type GormArticleDao struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDAO {
	return &GormArticleDao{
		db: db,
	}
}
func (dao *GormArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}
func (dao *GormArticleDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	//依赖gorm 忽略零值的特性,会根据主键进行跟新 ,可读性很差
	res := dao.db.WithContext(ctx).Model(&art).
		//这样的做法节省了数据库的查询 author_id
		Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		})
	//你要不要检查真的更新了没有

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		//dangerousDBOp.Count(1)
		//补充一些日志更加完整
		return fmt.Errorf("更新失败,可能是创作者非法 id:%d,author_id:%d", art.Id, art.AuthorId)
	}
	return res.Error

}

// Article 这是制作库的
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//长度
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	//TODO如何设置索引,在帖子里什么样的查询场景,对于创作者来说,是不是看草稿箱,看到所有自己的文章?产品经理告诉你,要按照创建的时间的倒叙排序
	//- 在 authorId  和 ctime上创建联合索引
	//在authorId 上创建索引
	AuthorId int64 `gorm:"index=aid_ctime"` //创建联合索引 index=aid_ctime
	Ctime    int64 `gorm:"idex=aid_ctime"`
	Utime    int64
	//TODO最佳选择就是在author_Id 和 Ctime联合创建联合索引
}
