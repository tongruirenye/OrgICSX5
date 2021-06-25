package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() {
	mysqluser, uerr := web.AppConfig.String("mysqluser")
	if uerr != nil {
		panic(uerr)
	}

	mysqlpass, perr := web.AppConfig.String("mysqlpass")
	if perr != nil {
		panic(perr)
	}

	mysqlurl, lerr := web.AppConfig.String("mysqlurl")
	if lerr != nil {
		panic(lerr)
	}

	mysqldb, derr := web.AppConfig.String("mysqldb")
	if derr != nil {
		panic(derr)
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&Asia%2FShanghai",
		mysqluser,
		mysqlpass,
		mysqlurl,
		mysqldb)
	if err := orm.RegisterDataBase("default", "mysql", url); err != nil {
		panic(err)
	}
	orm.RegisterModel(new(User), new(Sub))
	if err := orm.RunSyncdb("default", false, true); err != nil {
		panic(err)
	}

	NewAdmin()
}
