package framework

import (
	"fmt"
	"net/http"
	"strings"
)

// Core 框架核心结构
type Core struct {
	router      map[string]*Tree
	middlewares []ControllerHandler
}

//func treePrint(node *node, level int) {
//	println(node.segment, level)
//	if node.children != nil && len(node.children) > 0 {
//		for _, cnode := range node.children {
//			treePrint(cnode, level+1)
//		}
//	}
//}

//func (c *Core) PrintRoute() {
//	root := c.router["GET"].root
//	treePrint(root, 0)
//}

func (c *Core) Printmiddlewares() {
	print(len(c.middlewares))
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router, middlewares: []ControllerHandler{}}
}

//大小写不敏感

func (c *Core) GET(url string, handlers ...ControllerHandler) {

	allHandlers := append(c.middlewares, handlers...)

	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		//log.Fatal("add router error", err.Error())
		fmt.Println(err.Error())
	}
}

func (c *Core) POST(url string, handlers ...ControllerHandler) {

	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		//log.Fatal("add router error", err.Error())
		fmt.Println(err.Error())
	}
}

func (c *Core) PUT(url string, handlers ...ControllerHandler) {

	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		//log.Fatal("add router error", err.Error())
		fmt.Println(err.Error())
	}
}

func (c *Core) DELETE(url string, handlers ...ControllerHandler) {

	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		//log.Fatal("add router error", err.Error())
		fmt.Println(err.Error())
	}
}

func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Group 从core中初始化Group  ?? 但是这样只能构建一级group
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// FindRouteByRequest 匹配路由,没有匹配返回nil
func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
	//uri,method 全部转化为大写
	uri := request.URL.Path
	method := request.Method

	upperMethod := strings.ToUpper(method)

	//查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {

		return methodHandlers.FindHandler(uri)
	}
	return nil
}

func (c *Core) FindRouteNodeByRequest(request *http.Request) *node {
	uri := request.URL.Path
	method := request.Method

	upperMethod := strings.ToUpper(method)
	if routerTree, ok := c.router[upperMethod]; ok {
		matchNode := routerTree.root.matchNode(uri)
		return matchNode
	}
	return nil

}

//  框架核心结构实现Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	ctx := NewContext(request, response)

	matchNode := c.FindRouteNodeByRequest(request)

	handlers := c.FindRouteByRequest(request)
	if handlers == nil {
		ctx.SetOkStatus().Json("not found")
		return
	}
	params := matchNode.parseParamsFromEndNode(request.URL.Path)

	ctx.SetHandlers(handlers)
	ctx.SetParams(params)

	if err := ctx.Next(); err != nil {
		ctx.SetOkStatus().Json("inner error")
		return

	}

}
