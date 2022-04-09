package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/herrhu97/zinx/ziface"
)

type GlobalObj struct {
	 /*
        Server
    */
    TcpServer ziface.IServer //当前Zinx的全局Server对象
    Host      string         //当前服务器主机IP
    TcpPort   int            //当前服务器主机监听端口号
    Name      string         //当前服务器名称
    /*
        Zinx
    */
    Version          string //当前Zinx版本号
    MaxPacketSize    uint32 //都需数据包的最大值
    MaxConn          int    //当前服务器主机允许的最大链接个数
    WorkerPoolSize   uint32 //业务工作Worker池的数量
    MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
    /*
        config file path
    */
    ConfFilePath string
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
	//初始化GlobalObject变量，设置一些默认值
    GlobalObject = &GlobalObj{
        Name:          "ZinxServerApp",
        Version:       "V0.4",
        TcpPort:       7777,
        Host:          "0.0.0.0",
        MaxConn:       12000,
        MaxPacketSize: 4096,
        ConfFilePath:  "conf/zinx.json",
        WorkerPoolSize: 10,
        MaxWorkerTaskLen: 1024,
    }
    //从配置文件中加载一些用户配置的参数
    GlobalObject.Reload()

}
