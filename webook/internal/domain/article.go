package domain

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
