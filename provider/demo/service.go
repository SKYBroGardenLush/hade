package demo

import (
	"fmt"
	"github.com/SKYBroGardenLush/skycraper/framework"
)

type DemoService struct {
	//实现接口
	Service
	//参数
	c framework.Container
}

func (d *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}

//初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	fmt.Println("new demo services")
	return &DemoService{c: c}, nil
}
