package html

import (
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/models"
	"go-iptv/until"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admins(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}

	var pageData = dto.AdminsDto{
		LoginUser: username,
		Title:     "管理员设置",
	}

	dao.DB.Model(&models.IptvAdmin{}).Where("id = 1").First(&pageData.Admins)

	c.HTML(http.StatusOK, "admin_admins.html", pageData)
}

func Updata(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}

	var pageData = dto.UdataDto{
		LoginUser: username,
		Title:     "在线升级",
		Version:   until.GetVersion(),
	}

	c.HTML(http.StatusOK, "admin_updata.html", pageData)
}
