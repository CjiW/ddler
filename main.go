package main

import (
	"main/consts"
	"main/router"
	"main/sql"
	"main/tools"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	
	tools.InitDB()
	c := cron.New()
	_, _ = c.AddFunc(consts.AUTOTIME, sql.AutoRemind)
	c.Start()
	r := gin.Default()
	r.Use(tools.Cors())
	router.UseRouter(r)
	
	_ = r.Run(":10086")
	//tools.SendMsg(tools.RemindMsg("ou_14661677e00831a1b66d17f03e188706","ddler","2021/10/31",2),"card",[]string{"ou_14661677e00831a1b66d17f03e188706"})
}
