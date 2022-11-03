package net

//WsMsgResp =======返回数据结构体
type WsMsgResp struct {
	Body *RespBody
}

type RespBody struct {
	Seq  int64       `json:"seq"`  //序号
	Name string      `json:"name"` //名称
	Code int         `json:"code"` //状态码
	Msg  interface{} `json:"msg"`  //具体的消息
}

//WsMsgReq =====请求数据结构体
type WsMsgReq struct {
	Body *ReqBody
	Conn WSConnIface
}

type ReqBody struct {
	Seq   int64       `json:"seq"`  //序号
	Name  string      `json:"name"` //名称
	Msg   interface{} `json:"msg"`  //具体的消息
	Proxy string      `json:"proxy"`
}

//===========心跳机制

const HeartbeatMsg = "heartbeat"

type Heartbeat struct {
	CTime int64 `json:"ctime"`
	STime int64 `json:"stime"`
}

//Handshake 重新与客户端交流

const HandshakeMsg = "handshake"

type Handshake struct {
	Key string `json:"key"`
}
