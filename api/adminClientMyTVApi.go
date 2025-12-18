package api

import (
	"go-iptv/dto"
	"go-iptv/service"
	"go-iptv/until"

	"github.com/gin-gonic/gin"
)

func ClientMyTV(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.Request.ParseForm()
	params := c.Request.PostForm
	var res dto.ReturnJsonDto

	for k := range params {
		switch k {
		case "buildMyTV":
			res = service.SetMyTVAppInfo(params)
		}
	}
	c.JSON(200, res)
}

func BuildMyTVStatus(c *gin.Context) {
	_, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	c.JSON(200, service.GetMyTVBuildStatus())
}
