package http

import (
  "github.com/SKYBroGardenLush/skyscraper/app/http/module/demo"

  "github.com/SKYBroGardenLush/skyscraper/framework/contract"
  "github.com/SKYBroGardenLush/skyscraper/framework/gin"
  ginSwagger "github.com/SKYBroGardenLush/skyscraper/framework/middlewares/gin-swagger"
  "github.com/SKYBroGardenLush/skyscraper/framework/middlewares/static"
  swaggerFiles "github.com/swaggo/files"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

  // /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
  r.Use(static.Serve("/", static.LocalFile("./dist", false)))
  container := r.GetContainer()
  configService := container.MustMake(contract.ConfigKey).(contract.Config)

  if configService.GetBool("app.swagger") == true {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
  }
  // 动态路由定义
  demo.Register(r)
}
