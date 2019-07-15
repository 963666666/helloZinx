package znet

import (
	"fmt"
	"helloZinx/ziface"
	"net"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	msgHandler  ziface.IMsgHandler
	connManager ziface.IConnManager
	OnConnStart func(connection ziface.IConnection)
	OnConnStop  func(connection ziface.IConnection)
}

func (s *Server) Start() {

}

func (s *Server) Stop() {

}

func (s *Server) Server() {
	fmt.Printf("[START] Server name: %s,listener at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	go func() {
		//0 启动Worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " success, now listening...")

		//3 启动server网络连接业务

		var connId uint32
		connId = 0
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			if s.connManager.Len() >= 10 {
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, connId, s.msgHandler)
			connId++

			go dealConn.Start()
			/*buf := make([]byte, 512)
			cnt, err1 := conn.Read(buf)
			if err1 != nil {
				fmt.Println("recv buf err", err1)
				continue
			}

			fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)

			writeCnt, err3 := conn.Write(buf)
			if err3 != nil {
				fmt.Println("write buf err", err3)
			}
			fmt.Printf("send client buf %s, cnt %d\n", buf, writeCnt)*/
		}
	}()
}

/*
初始化
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        name,
		IPVersion:   "tcp4",
		IP:          "0.0.0.0",
		Port:        8999,
		msgHandler:  NewMsgHandler(),
		connManager: NewConnManager(),
	}

	return s
}

//得到连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.connManager
}

//设置该Server的连接创建时的Hook函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接端开始的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接onConnStart Hook函数
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("callOnConnStart")
		s.OnConnStart(connection)
	}
}

//调用连接onConnStop Hook函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("callOnConnStart")
		s.OnConnStop(connection)
	}
}

func init() {
	fmt.Println("helloZinx is running")
	fmt.Println("helloZinx is running")
	fmt.Println("helloZinx is running")
	fmt.Println("helloZinx is running")
	fmt.Println("helloZinx is running")
	fmt.Println("helloZinx is running")
}