package main

import (
  "github.com/SKYBroGardenLush/skycraper/framework/gin"
  "github.com/SKYBroGardenLush/skycraper/provider/demo"
)

// 注册路由规则
func registerRouterGin(core *gin.Engine) {
  // 需求1+2:HTTP方法+静态路由匹配
  //core.Use(middlewares.Recovery())
  core.GET("/user/login", GinUserLoginController)

  //core.PrintRoute()
  // 需求3:批量通用前缀
  subjectApi := core.Group("/subject")

  {
    // 需求4:动态路由
    //subjectApi.DELETE("/:id", SubjectDelController)
    //subjectApi.PUT("/:id", SubjectUpdateController)
    subjectApi.GET("/:id", GinSubjectGetController)
    //subjectApi.GET("/list/all", SubjectListController)
  }
  //core.PrintRoute()
}

func GinUserLoginController(c *gin.Context) {
  // 获取demo 服务实例

  demoService := c.MustMake(demo.Key).(*demo.DemoService)
  foo := demoService.GetFoo()

  c.ISetOkStatus().IJson(foo)

}

func GinSubjectDelController(c *gin.Context) {
  c.ISetOkStatus().IJson("ok, SubjectDelController")

}
func GinSubjectUpdateController(c *gin.Context) {
  c.ISetOkStatus().IJson("ok, SubjectUpdateController")

}
func GinSubjectGetController(c *gin.Context) {
  c.ISetOkStatus().IJson("ok, SubjectGetController")

}
func GinSubjectListController(c *gin.Context) {
  c.ISetOkStatus().IJson("ok, SubjectListController")

}
