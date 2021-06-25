package storage

import (
	"github.com/beego/beego/v2/server/web"
)

var DefaultStorage *WebDevClient

func InitStorage() {
	root, uerr := web.AppConfig.String("webdevroot")
	if uerr != nil {
		panic(uerr)
	}

	user, perr := web.AppConfig.String("webdevuser")
	if perr != nil {
		panic(perr)
	}

	pass, lerr := web.AppConfig.String("webdevpass")
	if lerr != nil {
		panic(lerr)
	}

	DefaultStorage = NewWebDevClient(root, user, pass)
	if DefaultStorage == nil {
		panic("storage error")
	}
}
