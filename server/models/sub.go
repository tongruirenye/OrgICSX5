package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Sub struct {
	Id      int       `orm:"auto;pk"`
	Name    string    `orm:"size(64):unique"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
}

func SubGetList() ([]*Sub, int64) {
	query := orm.NewOrm().QueryTable("sub")
	count, _ := query.Count()
	list := make([]*Sub, 0)
	query.All(&list)
	return list, count
}

func SubAdd(s *Sub) (int64, error) {
	return orm.NewOrm().Insert(s)
}

func SubGet(name string) (*Sub, error) {
	sub := Sub{Name: name}
	err := orm.NewOrm().Read(&sub, "Name")
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func NewTestSub() {
	if _, err := SubGet("读书.org"); err != nil {
		if err == orm.ErrNoRows {
			sub := &Sub{
				Name: "读书.org",
			}
			if _, e := SubAdd(sub); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
	if _, err := SubGet("考试.org"); err != nil {
		if err == orm.ErrNoRows {
			sub := &Sub{
				Name: "考试.org",
			}
			if _, e := SubAdd(sub); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
	if _, err := SubGet("工作.org"); err != nil {
		if err == orm.ErrNoRows {
			sub := &Sub{
				Name: "工作.org",
			}
			if _, e := SubAdd(sub); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
	if _, err := SubGet("项目.org"); err != nil {
		if err == orm.ErrNoRows {
			sub := &Sub{
				Name: "项目.org",
			}
			if _, e := SubAdd(sub); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
	if _, err := SubGet("运动.org"); err != nil {
		if err == orm.ErrNoRows {
			sub := &Sub{
				Name: "运动.org",
			}
			if _, e := SubAdd(sub); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
	if _, err := SubGet("日程.org"); err != nil {
		if err == orm.ErrNoRows {
			sub := &Sub{
				Name: "日程.org",
			}
			if _, e := SubAdd(sub); e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
}
