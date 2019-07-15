package znet

import (
	"fmt"
	"helloZinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connManager *ConnManager) Len() int {
	return len(connManager.connections)
}

func (connManager *ConnManager) Add(conn ziface.IConnection) {
	connManager.connLock.Lock()
	connManager.connections[conn.GetConnId()] = conn
	defer connManager.connLock.Unlock()
}

func (connManager *ConnManager) Remove(conn ziface.IConnection) {
	connManager.connLock.Lock()
	defer connManager.connLock.Unlock()
	delete(connManager.connections, conn.GetConnId())

	fmt.Println("connection remove connId =", conn.GetConnId(), "success connManager num is ", connManager.Len())
}