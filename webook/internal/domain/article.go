package domain

// TODO  可以同时表达制作库和线上库的概念吗,可以同时表达作者眼中的article 和读者眼中的article
type Article struct {
	Title   string
	Content string
	Author  Author
	Id      int64
}

type Author struct {
	Id   int64
	Name string
}
