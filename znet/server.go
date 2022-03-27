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
	Router    ziface.IRouter
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
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
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			panic(err)
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

			dealConn := NewConnection(conn, cid, s.Router)
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
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}
