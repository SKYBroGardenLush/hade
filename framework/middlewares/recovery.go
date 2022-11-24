package middlewares

import (
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		fmt.Println("hhhh")
		defer func() {
			if err := recover(); err != nil {
				c.SetOkStatus().Json(err)
			}
		}()

		c.Next()
		return nil
	}
}
