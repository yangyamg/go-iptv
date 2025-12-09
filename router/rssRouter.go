package router

import (
	"go-iptv/api"

	"github.com/gin-gonic/gin"
)

func RssRouter(r *gin.Engine, path string) {
	router := r.Group(path)
	{
		router.GET("/getRss/:token/paylist.m3u", api.GetTXTRssM3u)
		router.GET("/getRss/:token/paylist.txt", api.GetTXTRssTxt)
		router.GET("/ku9/:token/paylist.txt", api.GetTXTRssTxtKu9)
		router.GET("/epg/:token/e.xml", api.GetTXTRssEpg)
	}
}
