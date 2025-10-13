package web

import "Book_Exp/webook/internal/domain"

//对标前端

type RewardReq struct {
	Id  int64 `json:"id"`
	Amt int64 `json:"amt"`
}

// TODO 点赞和取消点赞准备复用这个
type LikeReq struct {
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
}

type CollectReq struct {
	Id  int64 `json:"id"`
	Cid int64 `json:"cid"`
}

type ArtcleVO struct {
	Title string `josn:"title"`
	//内容
	Content string `json:"content"`
	//摘要
	Abstract string `json:"abstract"`
	Author   string `json:"author"`
	Id       int64  `json:"id"`
	//状态这个可以是前端,也可以是后端来处理
	Status uint8  `json:"status"`
	Ctime  string `json:"ctime"`
	Utime  string `json:"utime"`
	//计数
	LikeCnt    int64 `json:"like_cnt"`
	ReadCnt    int64 `json:"read_cnt"`
	CollentCnt int64 `json:"collent_cnt"`
	//我个人有没有点过赞,和收藏
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ListReq struct {
	Offset int `json:"offest"`
	Limit  int `json:"limit"`
}
type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `josn:"title"`
	Content string `json:"content"`
}

func (rep ArticleReq) toDomain(uid int64) domain.Article {
	var req ArticleReq
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}

}
