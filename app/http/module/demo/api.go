package demo

import (
	demoService "github.com/SKYBroGardenLush/skycraper/app/provider/demo"
	"github.com/SKYBroGardenLush/skycraper/framework/contract"
	"github.com/SKYBroGardenLush/skycraper/framework/gin"
)

func Register(r *gin.Engine) error {
	r.Bind(&demoService.DemoProvider{})

	r.GET("/demo/demo", handlerDemo)

	return nil
}

func handlerDemo(c *gin.Context) {

	// 获取password
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	password := configService.GetString("database.mysql.password")
	// 打印出来
	c.JSON(200, password)

}
