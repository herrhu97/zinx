package znet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/herrhu97/zinx/utils"
	"github.com/herrhu97/zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// Router    ziface.IRouter
	msgHandler ziface.IMsgHandle
	ConnMgr    ziface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
	// =======================
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf error ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] ServerName: %s, Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	go func() {
		s.msgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}
		var cid uint32
		cid = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("listener accept tcp err ", err)
				continue
			}

			if s.ConnMgr.Len() > utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()

		}
	}()

}

func (s *Server) Serve() {
	s.Start()

	for {
		time.Sleep(10 * time.Minute)
	}
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
