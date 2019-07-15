package znet

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"helloZinx/ziface"
	"io"
	"net"
	"zinx/examples/zinx_release/zinx/znet"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer ziface.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//当前连接的ID 也可以称作SessionId， ID全局唯一
	ConnId uint32

	//当前连接的关闭状态
	isClose bool

	//当前连接所绑定的处理业务方法API
	msgHandler ziface.IMsgHandler

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool

	//无缓冲管道，用于读写两个Goroutine之间的消息通信
	msgChan chan []byte

	//有缓冲管道，用户读写两个Goroutine之间的消息通信
	msgBuffChan chan []byte
}

//初始化连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnId:       connId,
		isClose:      false,
		msgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, 100),
	}

	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")

	defer fmt.Println("connId = ", c.ConnId, "Reader is exit , remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		//创建拆包解包的对象
		dp := NewDataPack()
		//读取客户端的数据到buf中，最大512字节
		headData := make([]byte, dp.GetDataLen())

		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Println("recv buf err: ", err)
			break
		}
		//拆包得到msgId，dataLen
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unPack err ", err)
			break
		}
		//根据dataLen得到读取data，放到msg.data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("read msg data err ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		request := &Request{
			conn: c,
			msg:  msg,
		}

		//当前没有配置文件，直接交给连接池worker pool处理
		c.msgHandler.SendMsgToTaskQueue(request)
	}
}

func (c *Connection) StartWrite() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err", err, "conn writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send buff data err ", err, "conn writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn start ... ConnId = ", c.ConnId)

	go c.StartReader()

	go c.StartWrite()

}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	fmt.Println("Conn stop .. ConnId = ", c.ConnId)

	if c.isClose == true {
		return
	}
	c.isClose = true

	c.Conn.Close()

	c.ExitBuffChan <- true

	c.TcpServer.GetConnMgr().Remove(c)

	close(c.msgChan)
	close(c.msgBuffChan)
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
	if c.isClose == true {
		return errors.New("connection is closed when send msg")
	}
	//将data封包，发送
	dp := NewDataPack()
	msg, err := dp.Pack(znet.NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack err msgId is ", msgID)
		return err
	}

	//写回给客户端
	c.msgChan <- msg

	return nil
}
