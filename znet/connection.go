package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/herrhu97/zinx/utils"
	"github.com/herrhu97/zinx/ziface"
)

type Connection struct {
	TCPServer ziface.IServer

	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	MsgHandler ziface.IMsgHandle

	ExitBuffChan chan bool

	msgChan chan []byte

	msgBuffChan chan []byte

	property map[string]interface{}

	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32,
	msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn read exit!")
	defer c.Stop()

	for {
		dp := NewDataPack()
		buf := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), buf); err != nil {
			fmt.Println("recv msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}

		msg, err := dp.Unpack(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()

	c.TCPServer.CallOnConnStart(c)

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	fmt.Println("connection is stopping...")
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.TCPServer.CallOnConnStop(c)

	c.Conn.Close()
	c.ExitBuffChan <- true

	c.TCPServer.GetConnMgr().Remove(c)

	close(c.msgChan)
	close(c.ExitBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgBuffChan <- msg
	return nil
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}

		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
