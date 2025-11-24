package connpool

import (
	"context"
	"database/sql"
	"errors"
)

var errUnkownPattern = errors.New("未知的双写错误")

const (
	patternSrcOnly  = "SRC_ONLY"
	patternDstOnly  = "DST_ONLY"
	patternSrcFirst = "SRC_FIRST"
	patternDstFirst = "DST_FIRST"
)

type DoubleWirtePool struct {
}

func (d *DoubleWirtePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//TODO sql.Stmt 是一个结构体,你没有办法返回一个代表双写的 stmt
	panic("implement me")
}

func (d *DoubleWirtePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//增删改,或者非查询语句
	panic("implement me")
}

func (d *DoubleWirtePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWirtePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	//TODO implement me
	panic("implement me")
}
