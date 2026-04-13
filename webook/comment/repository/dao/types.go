package dao

import (
	"context"
)

//go:generate mockgen -source=./comment.go -package=daomocks -destination=mocks/comment.mock.go CommentDAO
type CommentDAO interface {
	Insert(ctx context.Context, u Comment) error
	// FindByBiz 只查找一级评论
	FindByBiz(ctx context.Context, biz string,
		bizId, minID, limit int64) ([]Comment, error)
	// FindCommentList Comment的id为0 获取一级评论，如果不为0获取对应的评论，和其评论的所有回复
	FindCommentList(ctx context.Context, u Comment) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]Comment, error)
	// Delete 删除本节点和其对应的子节点
	Delete(ctx context.Context, u Comment) error
	FindOneByIDs(ctx context.Context, id []int64) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
}
