package connpool

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDoubleWirtePool(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWirtePool{},
	}))
	require.NoError(t, err)
	t.Log(db)
}
