package router

/*
==============路由模块===================
*/

type Router struct {
	routerTree *tree         // all routers
	middleware []HandlerFunc //中间件函数不止一个
}

//NewRouter 路由函数构造
func NewRouter() *Router {
	return &Router{
		routerTree: newTree(),
		middleware: make([]HandlerFunc, 0),
	}
}

//Use 注册中间件
func (r *Router) Use(middlewares ...HandlerFunc) {
	r.middleware = append(r.middleware, middlewares...)
}

//AddRouter 添加 handlerFunc
func (r *Router) AddRouter(relativePath string, handlerFunc ...HandlerFunc) {
	allMiddlewares := append(r.middleware, handlerFunc...)
	//根据路径添加 handler
	r.routerTree.addRouter(relativePath, allMiddlewares...)
}

func (r *Router) Group(prefix string, middlewareFunc ...HandlerFunc) *Group {
	group := newGroup(r, prefix)
	group.middleware = middlewareFunc
	return group
}

func (r *Router) Run(c *Context) {
	//就是找到路由树中的handler 处理请求
	_ = r.routerTree.findHandler(c.req.Body.Name)

}
