package connpool

import (
	"testing"

	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDoubleWirtePool(t *testing.T) {
	webook, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	require.NoError(t, err)
	err = webook.AutoMigrate(&Interactive{})
	require.NoError(t, err)
	intr, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook_intr"))
	err = intr.AutoMigrate(&Interactive{})
	require.NoError(t, err)
	require.NoError(t, err)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWirtePool{
			src:     webook.ConnPool,
			dst:     intr.ConnPool,
			pattern: atomicx.NewValueOf(patternDstFirst),
		},
	}))
	require.NoError(t, err)
	t.Log(db)
	err = db.Create(&Interactive{
		Biz:   "test",
		BizId: 123,
	}).Error
	require.NoError(t, err)
	//TODO  事务问题
	err = db.Transaction(func(tx *gorm.DB) error {
		return db.Create(&Interactive{
			Biz:   "test_tx",
			BizId: 456,
		}).Error
	})
	require.NoError(t, err)
	t.Log(db)
	err = db.Model(&Interactive{}).Where("id > ?", 0).Updates(map[string]any{
		"biz_id": 789,
	}).Error
	require.NoError(t, err)

}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_id"`
	Biz        string `gomr:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	LikeCnt    int64
	Ctime      int64
	Utime      int64
}
