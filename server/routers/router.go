package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tongruirenye/OrgICSX5/server/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/sub", &controllers.SubController{})
}
