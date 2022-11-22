package framework

// IGroup 前缀分组接口
type IGroup interface {
	GET(string, ...ControllerHandler)
	POST(string, ...ControllerHandler)
	PUT(string, ...ControllerHandler)
	DELETE(string, ...ControllerHandler)

	// Group 实现嵌套group
	Group(string) IGroup

	Use(middlewares ...ControllerHandler)
}

// Group 实现IGroup
type Group struct {
	core        *Core
	parent      *Group
	prefix      string
	middlewares []ControllerHandler
}

// NewGroup 初始化group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		prefix:      prefix,
		parent:      nil,
		middlewares: []ControllerHandler{},
	}
}

func (g *Group) Group(uri string) IGroup {
	childGroup := NewGroup(g.core, uri)
	childGroup.parent = g
	return childGroup
}

//获取当前group的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

//获取当前group所有的中间件
func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) GET(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.middlewares, handlers...)
	g.core.GET(uri, allHandlers...)

}

func (g *Group) POST(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.middlewares, handlers...)
	g.core.POST(uri, allHandlers...)

}

func (g *Group) PUT(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.middlewares, handlers...)
	g.core.PUT(uri, allHandlers...)

}

func (g *Group) DELETE(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.middlewares, handlers...)
	g.core.GET(uri, allHandlers...)

}
