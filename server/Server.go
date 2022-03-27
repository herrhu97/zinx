package main

import "github.com/herrhu97/zinx/znet"

func main() {
	s := znet.NewServer("zinx 0.1")
	s.Serve()
}
