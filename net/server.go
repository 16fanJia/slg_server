package net

import (
	"github.com/gorilla/websocket"
	"net/http"
)

/*
=======封装 Server 服务 ============
*/

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Server net 链接服务
type Server struct {
	addr        string
	router      *Router
	needSecret  bool
	beforeClose func(WSConnIface)
}

func NewServer(addr string, needSecret bool) *Server {
	return &Server{
		addr:       addr,
		needSecret: needSecret,
	}
}

//SetBeforeCloseFunc 设置服务器的before close hook func
func (s *Server) SetBeforeCloseFunc(hookFunc func(WSConnIface)) {
	s.beforeClose = hookFunc
}

//RegisterRouter 注册路由函数
func (s *Server) RegisterRouter(router *Router) {
	s.router = router
}

//Start 开启server服务
func (s *Server) Start() (err error) {
	http.HandleFunc("/", s.wsHandler)
	err = http.ListenAndServe(s.addr, nil)
	return err
}

func (s *Server) wsHandler(resp http.ResponseWriter, req *http.Request) {
	//http升级为 websocket 连接
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	conn := NewServerConn(s, wsSocket, s.needSecret)
	conn.Start()

}
