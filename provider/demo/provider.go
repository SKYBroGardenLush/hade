package demo

import (
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework"
)

//服务提供方
type DemoServiceProvider struct {
}

// Name 将服务对应的字符串凭证返回
func (d *DemoServiceProvider) Name() string {
	return Key
}

// Register 注册初始化服务实例的方法
func (d *DemoServiceProvider) Register(c framework.Container) framework.NewInstance {
	return NewDemoService
}

// IsDefer 是否延迟实例化
func (d *DemoServiceProvider) IsDefer() bool {
	return true
}

// Params 表示实例化的参数
func (d *DemoServiceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (d *DemoServiceProvider) Boot(c framework.Container) error {
	fmt.Println("demo services boot")
	return nil
}
