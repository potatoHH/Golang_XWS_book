package connpool

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

//TODO 读写分离的功能

var errUnkonwPattern = errors.New("读写分离模式错误")

type WirteSplit struct {
	master gorm.ConnPool
	slaves []gorm.ConnPool
}

func (w *WirteSplit) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	tx, err := w.master.(gorm.TxBeginner).BeginTx(ctx, opts)
	if err != nil {

	}
}

func (w *WirteSplit) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//可以默认返回 master 或者 默认返回 salves
	return w.master.PrepareContext(ctx, query)
}

func (w *WirteSplit) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//默认走master
	return w.master.ExecContext(ctx, query, args...)
}

func (w *WirteSplit) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	//salves 要考虑 负载均衡，可以轮询的方式
	//TODO 这边可以玩骚操作，轮询，加权轮询，冰花的加权轮询，随机，加权随机，动态判定slaves 健康情况的均衡负载策略（永远挑选反应时间最快的 slave,或者暂时禁用掉已超时的slave  ）
}

func (w *WirteSplit) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	//TODO implement me
	panic("implement me")
}
