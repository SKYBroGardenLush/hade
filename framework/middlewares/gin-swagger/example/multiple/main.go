package main

import (
	"github.com/SKYBroGardenLush/skyscraper/framework/gin"
	ginSwagger "github.com/SKYBroGardenLush/skyscraper/framework/middlewares/gin-swagger"
	v1 "github.com/SKYBroGardenLush/skyscraper/framework/middlewares/gin-swagger/example/multiple/api/v1"
	v2 "github.com/SKYBroGardenLush/skyscraper/framework/middlewares/gin-swagger/example/multiple/api/v2"
	_ "github.com/SKYBroGardenLush/skyscraper/framework/middlewares/gin-swagger/example/multiple/docs"
	swaggerFiles "github.com/swaggo/files"
)

func main() {
	// New gin router
	router := gin.New()

	// Register api/v1 endpoints
	v1.Register(router)
	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	// Register api/v2 endpoints
	v2.Register(router)
	router.GET("/swagger/v2/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v2")))

	// Listen and Server in
	_ = router.Run()
}
