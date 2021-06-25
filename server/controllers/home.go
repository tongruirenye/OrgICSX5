package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Data["username"] = c.GetSession("username")
	c.TplName = "site.tpl"
}
