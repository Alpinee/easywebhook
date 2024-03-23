/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router = gin.Default()

func InitRouter() {
	// Router.HTMLRender = utils.LoadTemplateFiles("templates", ".html")
	Router.LoadHTMLGlob("./html/**/*")
	Router.Use(AuthMiddleware([]string{"/login", "/api/login", "/webhook"}))

	// 注册页面
	{
		Router.GET("/", Index)

		Router.GET("/login", func(context *gin.Context) {
			context.HTML(http.StatusOK, "login", gin.H{})
		})
	}

	Router.Any("/webhook", handleWebhook)
	// 注册api
	api := Router.Group("/api")
	api.POST("/login", login)
	{
		script := api.Group("/script")
		script.POST("/add", addScript)
		script.GET("/get/:id", getScript)
		script.PUT("/update/:id", updateScript)
		script.DELETE("/delete/:id", deleteScript)
	}

	// 指定监听的端口，例如8888
	port := ":8888"

	// 启动服务，监听指定端口
	err := Router.Run(port)
	if err != nil {
		panic("Start server failed.Reason:" + err.Error())
	}
}
