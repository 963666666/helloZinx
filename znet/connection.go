package znet

import (
	"fmt"
	"hello_zinx/ziface"
	"net"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//当前连接的ID 也可以称作SessionId， ID全局唯一
	ConnId uint32

	//当前连接的关闭状态
	isClose bool

	//当前连接所绑定的处理业务方法API
	handleApi ziface.HandelFunc

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}

//初始化连接的方法
func NewConnection(conn *net.TCPConn, connId uint32, callBackApi ziface.HandelFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnId:       connId,
		handleApi:    callBackApi,
		isClose:      false,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")

	defer fmt.Println("connId = ", c.ConnId, "Reader is exit , remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)

		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err: ", err)
			continue
		}
		//调用当前连接所绑定的HandleAPI
		if err := c.handleApi(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnId", c.ConnId, "handle is error:", err)
			break
		}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn start ... ConnId = ", c.ConnId)

	go c.StartReader()

}
//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	fmt.Println("Conn stop .. ConnId = ", c.ConnId)

	if c.isClose == true {
		return
	}
	c.isClose = true

	c.Conn.Close()

	close(c.ExitBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接的ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message 发送给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	return nil
}
