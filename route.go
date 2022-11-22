package main

import (
	"github.com/SKYBroGardenLush/skycraper/framework"
	"github.com/SKYBroGardenLush/skycraper/framework/middlewares"
	"time"
)

// 注册路由规则
func registerRouter(core *framework.Core) {
	// 需求1+2:HTTP方法+静态路由匹配
	//core.Use(middlewares.Recovery())
	core.GET("/user/login", middlewares.Recovery(), UserLoginController)

	//core.PrintRoute()
	// 需求3:批量通用前缀
	subjectApi := core.Group("/subject")
	subjectApi.Use(middlewares.Recovery())
	{
		// 需求4:动态路由
		//subjectApi.DELETE("/:id", SubjectDelController)
		//subjectApi.PUT("/:id", SubjectUpdateController)
		subjectApi.GET("/:id", SubjectGetController)
		//subjectApi.GET("/list/all", SubjectListController)
	}
	//core.PrintRoute()
}

func UserLoginController(c *framework.Context) error {
	time.Sleep(10 * time.Second)
	c.SetOkStatus().Json("ok, UserLoginController")
	return nil
}

func SubjectDelController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectDelController")
	return nil
}
func SubjectUpdateController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectUpdateController")
	return nil
}
func SubjectGetController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectGetController")
	return nil
}
func SubjectListController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectListController")
	return nil
}
