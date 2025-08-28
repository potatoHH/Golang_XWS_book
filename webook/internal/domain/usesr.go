package domain

import "time"

//user领域对象

type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Nickname string
	AboutMe  string
	Birthday time.Time
	Ctime    time.Time
	//Unionid  string
	//Openid   string
	//不要合并,万一以后有同样的字段名
	WechatInfo WecahteInfo
}

func (u *User) NewUser() {

}
