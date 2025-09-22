package article

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	UpSert(ctx context.Context, art PublishArticle) error
	SyncStatus(ctx context.Context, id int64, author int64, status uint8) error
	GetByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]Article, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (Article, error)
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
}
type GormArticleDao struct {
	db *gorm.DB
}

func (dao *GormArticleDao) GetPubById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GormArticleDao) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GormArticleDao) GetByID(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GormArticleDao) GetByAuthor(ctx context.Context, author int64, limit int, offset int) ([]Article, error) {
	var arts []Article
	// 在设计order by  语句的时候,要注意order by 中命中数据命中索引
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("author_id=?", author).
		Offset(offset).Limit(limit).
		//Order("utime DESC").
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "utime"}, Desc: true},
			//{Column: clause.Column{Name: "ctime"}, Desc: false},
		}}).
		Find(&arts).Error
	return arts, err

}

func (dao *GormArticleDao) SyncStatus(ctx context.Context, id int64, author int64, status uint8) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? AND author_id = ?", id, author).
			Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}

		res = tx.Model(&PublishedArticle{}).
			Where("id=? AND author_id = ?", id, author).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		return nil
	})
}

// UpSert nsert OR Update
func (dao *GormArticleDao) UpSert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	//这个是插入, clauses.Onconfilict 是数据冲突了
	err := dao.db.Clauses(clause.OnConflict{
		// 那些列冲突
		//Columns:[]clause.Column{
		//	clause.Column{Name: "id"},
		//},
		//数据冲突,啥也不干  DoNothing:
		//数据冲突了,并且符合where 条件的就会执行更新 DoUpdate  Where:
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		}),
	}).Create(&art).Error
	// MySql 最终的语句 Insert xxx On Duplicate key Update xxx   如果发生数据冲突就执行更新
	// 一条 msyql 不需要开事务
	// auto commit 自动提交
	return err
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
func (dao *GormArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	//TODO 先操作制作库(制作表),在操作线上库(线上表)
	//在事务内部采用闭包的形态
	//在gorm 帮助我们管理了事务的声明
	var id = art.Id
	//tx=> transaction
	//TODO begin,commit,rollback 都不需要我们去管理
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewGORMArticleDAO(tx)
		if art.Id > 0 {
			err = txDAO.UpdateById(ctx, art)

		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		//操作线上库
		return txDAO.UpSert(ctx, PublishArticle{Article: art})
	})
	return id, err

}

// Article 这是制作
