package ics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tongruirenye/OrgICSX5/server/config"
)

type GenIcsForm struct {
	Sign string `form:"sign" binding:"required"`
}

func GenIcs(c *gin.Context) {
	var form GenIcsForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusOK, "caonimade,tingdedongba")
		return
	}

	if form.Sign != config.AppConfig.Sign {
		c.String(http.StatusOK, "nidayede, tingdedongba")
		return
	}

	icsParser := c.MustGet("ics").(*ICS)
	icsParser.Task()
	c.String(http.StatusOK, "I love huhong!!!")
}

func UseIcs(ics *ICS) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ics", ics)
		c.Next()
	}
}
