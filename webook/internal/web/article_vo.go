package web

import "Book_Exp/webook/internal/domain"

//对标前端

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
