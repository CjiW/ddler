package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"main/consts"
	"main/router"
	"main/sql"
)

func main() {
	c := cron.New()
	_, _ = c.AddFunc(consts.AUTOTIME, sql.AutoRemind)
	c.Start()
	r := gin.Default()
	router.UseRouter(r)
	_ = r.Run(":10086")
	//tools.SendMsg(tools.RemindMsg("ou_14661677e00831a1b66d17f03e188706","ddler","2021/10/31",2),"card",[]string{"ou_14661677e00831a1b66d17f03e188706"})
}
