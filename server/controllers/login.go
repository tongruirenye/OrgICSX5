package controllers

import (
	"html/template"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tongruirenye/OrgICSX5/server/models"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	web.ReadFromRequest(&c.Controller)

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.tpl"
}

func (c *LoginController) Post() {
	username := c.GetString("username")
	password := c.GetString("password")

	errMsg := ""
	user, err := models.UserGet(username)
	if err != nil {
		errMsg = "用户不存在"
	} else if models.UserVerifyPassword(password, user.Salt, user.Password) == false {
		errMsg = "密码错误"
	}

	if errMsg == "" {
		c.SetSession("uid", user.Id)
		c.SetSession("username", user.Email)
		c.Redirect("/home", 302)
	} else {
		flash := beego.NewFlash()
		flash.Error(errMsg)
		flash.Store(&c.Controller)
		c.Redirect("/login", 302)
	}
}
