package html

import (
	"encoding/json"
	"go-iptv/dao"
	"go-iptv/dto"
	"go-iptv/until"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func License(c *gin.Context) {
	username, ok := until.GetAuthName(c)
	if !ok {
		c.JSON(200, dto.NewAdminRedirectDto())
		return
	}
	var pageData = dto.AdminLicenseDto{
		LoginUser: username,
		Title:     "进阶功能",
	}

	res, err := dao.WS.SendWS(dao.Request{Action: "reloadLic"})
	if err == nil {
		if err := json.Unmarshal(res.Data, &dao.Lic); err != nil {
			log.Println("license信息解析错误:", err)
		}
	}

	cfg := dao.GetConfig()
	pageData.Proxy = cfg.Proxy.Status
	pageData.ProxyAddr = cfg.Proxy.PAddr
	pageData.Lic = dao.Lic
	pageData.Port = cfg.Proxy.Port
	pageData.Lic.ExpStr = time.Unix(pageData.Lic.Exp, 0).Format("2006-01-02 15:04:05")

	c.HTML(200, "admin_license.html", pageData)
}
