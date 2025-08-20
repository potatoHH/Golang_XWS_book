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
}

func (u *User) NewUser() {

}
