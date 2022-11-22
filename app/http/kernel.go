package http

import "github.com/SKYBroGardenLush/skycraper/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	Routes(r)
	// 返回绑定路由后的Web引擎
	return r, nil
}
