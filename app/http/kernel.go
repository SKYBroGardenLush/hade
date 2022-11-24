package http

import (
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetContainer(container)
	Routes(r)
	// 返回绑定路由后的Web引擎
	return r, nil
}
