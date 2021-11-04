package router

import (
	"github.com/gin-gonic/gin"
	"main/router/api"
)

func UseRouter(r *gin.Engine) {
	r.POST("/msg", api.HandleMsg)
	// r.POST("/card", api.HandleCard)
	r.GET("/h5get", api.HandleH5get)
	r.POST("/h5post", api.HandleH5post)
}
