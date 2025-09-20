package article

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	// 标题的长度
	// 正常都不会超过这个长度
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 作者
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	Utime    int64 `bson:"utime,omitempty" gorm:"index"`
}
type ArticleV1 struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//长度
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	//TODO如何设置索引,在帖子里什么样的查询场景,对于创作者来说,是不是看草稿箱,看到所有自己的文章?产品经理告诉你,要按照创建的时间的倒叙排序
	//- 在 authorId  和 ctime上创建联合索引
	//在authorId 上创建索引
	AuthorId int64 `gorm:"index=aid_ctime"` //创建联合索引 index=aid_ctime
	Ctime    int64 `gorm:"index=aid_ctime"`
	Utime    int64
	//TODO最佳选择就是在author_Id 和 Ctime联合创建联合索引
	Status uint8
}

// PublishedArticle 衍生类型，偷个懒
type PublishedArticle Article
