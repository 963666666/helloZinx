package ziface

import "net"

type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态M
	Stop()

	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn

	//获取当前连接的ID
	GetConnId() uint32

	//获取远程客户端地址信息
	RemoteAddr() net.Addr

	//直接将Message 发送给远程的TCP客户端
	SendMsg(msgID uint32, data []byte) error

}

type HandelFunc func(*net.TCPConn, []byte, int) error
