package router

import "slg_server/net"

/*
 =======借用gin的思想 封装请求和返回==============
*/

type Context struct {
	req  *net.WsMsgReq
	resp *net.WsMsgResp

	handlers []HandlerFunc //当前请求的handler 链条
	idx      int           //当前请求调用到调用链的那个节点
}

//NewContext 初始化一个context
func NewContext(req *net.WsMsgReq, resp *net.WsMsgResp) *Context {
	return &Context{
		req:      req,
		resp:     resp,
		handlers: []HandlerFunc{},
		idx:      -1,
	}
}

//SetHandlers 为context 设置handlers
func (ctx *Context) SetHandlers(handlers []HandlerFunc) {
	ctx.handlers = handlers
}

//Next 调用context的下一个函数
func (ctx *Context) Next() {
	ctx.idx++
	if ctx.idx < len(ctx.handlers) {
		ctx.handlers[ctx.idx](ctx)
	}
}
