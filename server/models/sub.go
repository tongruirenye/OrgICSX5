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
