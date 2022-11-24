package app

import "github.com/SKYBroGardenLush/skyscraper/framework/gin"

// AppMiddleware 代表中间件函数
func AppMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		//您的处理代码
		context.Next()
	}
}
