package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"slg_server/util"
	"sync"
	"time"
)

/*
=======封装 ServerConn 连接服务 ============
*/

type ServerConn struct {
	server       *Server
	wsSocket     *websocket.Conn
	needSecret   bool
	property     map[string]interface{} //链接属性
	propertyLock sync.RWMutex           //保护链接属性读取和修改的锁
	isClosed     bool                   //链接是否关闭
	outChan      chan *WsMsgResp        //协程之间通信 写数据
	beforeClose  func(WSConnIface)      //关闭连接之前的函数
}

//NewServerConn 构建一个新的连接
func NewServerConn(s *Server, conn *websocket.Conn, needSecret bool) *ServerConn {
	return &ServerConn{
		server:     s,
		wsSocket:   conn,
		needSecret: needSecret,
		property:   make(map[string]interface{}),
		isClosed:   false,
		outChan:    make(chan *WsMsgResp, 1000), // buffed channel
	}
}

//SetProperty 设置连接属性
func (sc *ServerConn) SetProperty(key string, value interface{}) {
	sc.propertyLock.Lock()
	defer sc.propertyLock.Unlock()
	sc.property[key] = value
}

//GetPropertyByKey 根据key 获取对应的属性
func (sc *ServerConn) GetPropertyByKey(key string) (interface{}, error) {
	sc.propertyLock.RLock()
	defer sc.propertyLock.RUnlock()

	if val, ok := sc.property[key]; ok {
		return val, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//RemoveProperty 移除一个属性
func (sc *ServerConn) RemoveProperty(key string) {
	sc.propertyLock.Lock()
	defer sc.propertyLock.Unlock()

	delete(sc.property, key)
}

func (sc *ServerConn) GetAddr() string {
	return sc.wsSocket.RemoteAddr().String()
}

//Push 用于向通道中发送数据
func (sc *ServerConn) Push(name string, data interface{}) {
	resp := &WsMsgResp{Body: &RespBody{
		Seq:  0,
		Name: name,
		Msg:  data,
	}}
	sc.outChan <- resp
}

//Start 开启异步 读数据 和 写数据
func (sc *ServerConn) Start() {
	go sc.wsReadLoop()
	go sc.wsWriteLoop()
}

//wsReadLoop 循环读数据
func (sc *ServerConn) wsReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			//记日志
			fmt.Println(err)
			//关连接
			sc.Close()
		}
	}()
	for {
		//读一个么message
		_, msgData, err := sc.wsSocket.ReadMessage()
		if err != nil {
			break
		}
		data, err := util.UnZip(msgData)
		if err != nil {
			continue
		}

		//请求
		body := &ReqBody{}
		//检测是否加密
		if sc.needSecret {
			//获取秘钥
			if secretKey, err := sc.GetPropertyByKey("secretKey"); err == nil {
				key := secretKey.(string)
				//解密数据
				encrypt, err := util.AesCBCDecrypt(data, []byte(key), util.ZEROS_PADDING)
				if err != nil {
					sc.Handshake()
				}
				data = encrypt
			}
		} else {
			fmt.Println("未找到秘钥，连接需要重置，并且设置秘钥")
			sc.Handshake()
			return
		}

		//解析数据
		if err = json.Unmarshal(data, body); err == nil {
			req := &WsMsgReq{
				Body: body,
				Conn: sc,
			}

			resp := &WsMsgResp{Body: &RespBody{
				Seq:  body.Seq,
				Name: body.Name,
			}}
			//TODO 后续是否可以将req 和 resp 封装到 context 中

			//判断是否为心跳消息
			if body.Name == HeartbeatMsg {
				h := &Heartbeat{}
				//TODO 将 body中的msg字段 解析为h

				h.STime = time.Now().UnixNano() / 1e6
				resp.Body.Msg = h
			} else {
				//处理请求
				if sc.server.router != nil {
					sc.server.router.Run(req, resp)
				}
			}
			sc.outChan <- resp
		} else {
			fmt.Println("unmarshal error ", err)
			sc.Handshake()
		}
	}
	sc.Close()
}

//wsWriteLoop 循环写数据
func (sc *ServerConn) wsWriteLoop() {
	defer func() {
		if err := recover(); err != nil {
			sc.Close()
		}
	}()

	for {
		select {
		//取一个数据
		case msg := <-sc.outChan:
			//写给websocket
			sc.write(msg)
		}
	}
}

func (sc *ServerConn) write(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	if sc.needSecret {
		if secretKey, err := sc.GetPropertyByKey("secretKey"); err == nil {
			key := secretKey.(string)
			//加密数据
			encrypt, err := util.AesCBCEncrypt(data, []byte(key), util.ZEROS_PADDING)
			if err != nil {
				//记日志
			}
			data = encrypt
		}
	}

	if data, err = util.Zip(data); err != nil {
		//记日志
		return
	}
	//数据写入 websocket 中
	if err = sc.wsSocket.WriteMessage(websocket.BinaryMessage, data); err != nil {
		//记日志
		return
	}
}

//Close 关闭连接
func (sc *ServerConn) Close() {
	if !sc.isClosed {
		//执行关闭之前回调函数
		if sc.beforeClose != nil {
			sc.server.beforeClose(sc)
		}

		//关闭连接 回收资源
		sc.wsSocket.Close()
		sc.isClosed = true
	}
}

//Handshake 握手协议
func (sc *ServerConn) Handshake() {
	secretKey := ""
	if sc.needSecret {
		key, err := sc.GetPropertyByKey("secretKey")
		if err != nil {
			//生成一个随机的16位数
			secretKey = util.RandSeq(16)
		} else {
			secretKey = key.(string)
		}
	}

	handshake := &Handshake{Key: secretKey}
	body := &RespBody{
		Name: HandshakeMsg,
		Msg:  handshake,
	}
	var (
		data []byte
		err  error
	)

	if data, err = json.Marshal(body); err != nil {
		fmt.Println("handshake Marshal body error", err)
		return
	}

	if secretKey != "" {
		//需要加密
		sc.SetProperty("secretKey", secretKey)
	} else {
		//不需要加密
		sc.RemoveProperty("secretKey")
	}

	if data, err = util.Zip(data); err != nil {
		fmt.Println("zip data err", err)
		return
	}
	sc.wsSocket.WriteMessage(websocket.BinaryMessage, data)
}
