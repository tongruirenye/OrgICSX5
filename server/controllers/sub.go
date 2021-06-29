package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/task"
	"github.com/tongruirenye/OrgICSX5/server/ics"
	"github.com/tongruirenye/OrgICSX5/server/models"
	"github.com/tongruirenye/OrgICSX5/server/storage"
)

type SubController struct {
	beego.Controller
}

type SubList struct {
	Status string
	Name   string
}

func (c *SubController) Get() {
	files, err := storage.DefaultStorage.ListFileList("org/roam/project")
	if err != nil {
		c.Data["error"] = err.Error()
	} else {
		var subList []*SubList
		if files != nil {
			for _, file := range files {
				subList = append(subList, &SubList{Status: "订阅", Name: file})
			}
		}

		found := false
		lists, _ := models.SubGetList()
		if lists != nil && len(lists) > 0 {
			for _, sub := range lists {
				found = false
				if subList != nil {
					for _, sl := range subList {
						if sl.Name == sub.Name {
							sl.Status = "取消订阅"
							found = true
							break
						}
					}
				}
				if !found {
					subList = append(subList, &SubList{Status: "删除", Name: sub.Name})
				}
			}
		}

		c.Data["files"] = subList

	}

	c.Layout = "site.tpl"
	c.TplName = "sub.tpl"
}

type SubGenController struct {
	beego.Controller
}

func (c *SubGenController) Get() {
	task := task.NewTask("genics", "0 0 6 * * *", ics.GenIcsTask)
	task.Run(nil)
	c.Redirect("/sub", 302)
}
