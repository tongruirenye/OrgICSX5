package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/tongruirenye/OrgICSX5/server/utils"
)

type User struct {
	Id       int       `orm:"auto;pk"`
	Email    string    `orm:"size(64);unique"`
	Password string    `orm:"size(32)"`
	Salt     string    `orm:"size(6)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func NewAdmin() {
	if _, err := UserGet("tongruirenye@163.com"); err != nil {
		if err == orm.ErrNoRows {
			password, salt := UserNewPassword("tongruirenye")
			user := &User{
				Email:    "tongruirenye@163.com",
				Password: password,
				Salt:     salt,
			}
			if _, e := UserAdd(user); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
}

func UserAdd(u *User) (int64, error) {
	return orm.NewOrm().Insert(u)
}

func UserGet(email string) (*User, error) {
	user := User{Email: email}
	err := orm.NewOrm().Read(&user, "Email")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserNewPassword(password string) (string, string) {
	passwordMd5 := utils.Md5(password)
	salt := utils.GenerateRandomString(6)
	md5Password := utils.Md5(passwordMd5 + salt)
	return md5Password, salt
}

func UserVerifyPassword(password, salt, md5Password string) bool {
	passwordMd5 := utils.Md5(password)
	return md5Password == utils.Md5(passwordMd5+salt)
}
