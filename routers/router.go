package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(app *gin.Engine) {
	//首页
	app.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })

}
