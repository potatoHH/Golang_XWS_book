package domain

type Interactive struct {
	//计数
	LikeCnt    int64 `json:"like_cnt"`
	ReadCnt    int64 `json:"read_cnt"`
	CollectCnt int64 `json:"collent_cnt"`
	//我个人有没有点过赞,和收藏
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
	Biz       string
	BizId     int64
}
