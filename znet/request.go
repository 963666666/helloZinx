package znet

import (
	"helloZinx/ziface"
)

type Request struct {
	conn ziface.IConnection
	msg ziface.IMessage
}

func (request *Request)GetConnection() ziface.IConnection {
	return request.conn
}

func (request *Request)GetData() []byte {
	return request.msg.GetData()
}

func (request *Request) GetMsgId() uint32 {
	return request.msg.GetMsgId()
}
