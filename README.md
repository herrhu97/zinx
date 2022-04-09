# zinx

Simple Conccurency Tcp Server Develop Framework.

## Demo
`Shell`
~~~
mkdir conf && touch conf/zinx.json
echo "{
  "Name":"zinx v-0.10 demoApp",
  "Host":"127.0.0.1",
  "TcpPort":7777,
  "MaxConn":3,
  "WorkerPoolSize":10,
  "LogDir": "./mylog",
  "LogFile":"zinx.log"
}" > conf/zinx.json
~~~

`Server`
~~~
package main

import (
	"fmt"

	"github.com/herrhu97/zinx/ziface"
	"github.com/herrhu97/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()

	s.AddRouter(0, &PingRouter{})

	s.Serve()
}
~~~

`Client`
~~~
package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/herrhu97/zinx/znet"
)

func main() {

	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for n := 3; n >= 0; n-- {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx Client Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error")
			break
		}
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.MsgID, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
~~~
