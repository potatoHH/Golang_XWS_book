package domain

import "time"

//user领域对象

type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}

func (u *User) NewUser() {

}
