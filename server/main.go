package main

import (
	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tongruirenye/OrgICSX5/server/middleware"
	"github.com/tongruirenye/OrgICSX5/server/models"
	_ "github.com/tongruirenye/OrgICSX5/server/routers"
	"github.com/tongruirenye/OrgICSX5/server/storage"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.EnableXSRF = true
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.LoginVerify)
	models.InitDB()
	storage.InitStorage()
	go web.SetStaticPath("/static", "public")
	beego.Run()
}
