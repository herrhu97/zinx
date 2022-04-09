package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/herrhu97/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //map[connID]connection
	connLock    sync.RWMutex                  //读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection remove from ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	fmt.Println("clear all connections successfully")
}

