package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func LoginVerify(ctx *context.Context) {
	_, ok := ctx.Input.Session("uid").(int)
	if !ok && ctx.Request.RequestURI != "/login" {
		ctx.Redirect(302, "/login")
	}
}
