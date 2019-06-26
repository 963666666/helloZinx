package znet

import (
	"hello_zinx/ziface"
	ziface2 "zinx/ziface"
)

type Request struct {
	conn ziface.IConnection
	msg ziface2.IMessage
}
