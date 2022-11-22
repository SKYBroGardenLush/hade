package user

import "github.com/SKYBroGardenLush/skycraper/framework/gin"

// UserMiddleware 代表中间件函数
func UserMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
	}
}
