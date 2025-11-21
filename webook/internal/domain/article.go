package domain

import "time"

// TODO  可以同时表达制作库和线上库的概念吗,可以同时表达作者眼中的article 和读者眼中的article
type Article struct {
	Title   string
	Content string
	Author  Author
	Id      int64
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

const (
	//ArticleStatusUnknown  为了避免零值之类的问题
	ArticleStatusUnknown     = iota //位置状态
	ArticleStatusUnpublished        //未发表
	ArticleStatusPublished          //未发表
	ArticleStatusPrivate            // 仅自己可见
)

type ArticleStatus uint8

func (a Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	return string(cs[:100])
}

func (s ArticleStatus) ToUnit8() uint8 {
	return uint8(s)
}
func (s ArticleStatus) Vaild() bool { //判断是否合法
	return s.ToUnit8() > 0
}
func (s ArticleStatus) NonPublish() bool {
	return true
}
func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusUnpublished:
		return "unpublished"
	case ArticleStatusPublished:
		return "published"
	case ArticleStatusPrivate:
		return "private"
	default:
		return "unknow"
	}
}

// TODO  如果你的状态很复杂,有很多行为(搞了好多方法),转台里面需要一些些额外的字段,就用这个结构体
type ArticleStatusV1 struct {
	Val  uint8
	Name string
}

var ArticleStatusUnknow = ArticleStatusV1{Val: 0, Name: "unknow"}

type Author struct {
	Id   int64
	Name string
}
