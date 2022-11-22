package http

import (
	"github.com/SKYBroGardenLush/skycraper/app/http/module/demo"
	"github.com/SKYBroGardenLush/skycraper/framework/gin"
	"github.com/SKYBroGardenLush/skycraper/framework/middlewares/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	// 动态路由定义
	demo.Register(r)
}
