package router

type HandlerFunc func(c *Context)

type Group struct {
	router     *Router       //总router
	parent     *Group        //嵌套group的 父group
	prefix     string        //前缀
	middleware []HandlerFunc //中间件函数不止一个
}

func newGroup(r *Router, prefix string) *Group {
	return &Group{
		router:     r,
		parent:     nil,
		prefix:     prefix,
		middleware: []HandlerFunc{},
	}
}

//Use 注册中间件
func (g *Group) Use(middleware ...HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

//AddRouter 添加路由
func (g *Group) AddRouter(relativePath string, handlerFunc ...HandlerFunc) {
	//根据绝对的path 往路由树中添加路由
	path := g.getAbsolutePath() + relativePath
	//这个路径上的所有 middlewares
	allMiddlewares := append(g.getMiddlewares(), handlerFunc...)
	g.router.routerTree.addRouter(path, allMiddlewares...)
}

//Group 实现group中嵌套group
func (g *Group) Group(relativePath string) *Group {
	cGroup := newGroup(g.router, relativePath)
	cGroup.parent = g
	return cGroup
}

//getAbsolutePath 获取绝对的路径
/*
 g := r.Group("/name")
 g.AddRouter("/fanjia",..)
 path = /name/fanjia
*/
func (g *Group) getAbsolutePath() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePath() + g.prefix
}

//getMiddlewares 获取所有的中间件函数
func (g *Group) getMiddlewares() []HandlerFunc {
	if g.parent == nil {
		//最上层 router 的 middlewares
		return g.middleware
	}
	return append(g.parent.getMiddlewares(), g.middleware...)
}
