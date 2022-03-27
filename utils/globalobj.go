package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/herrhu97/zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer

	Host          string
	TcpPort       int
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       int
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("read conf file err ", err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		fmt.Println("json unmarshal conf data err ", err)
		os.Exit(1)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}
	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()

}
