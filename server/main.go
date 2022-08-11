package main

import (
	"flag"
	"p2p-test/server/tcp"
	"p2p-test/server/udp"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 13520, "server listen port")
	flag.Parse()
}

func main() {
	go tcp.RunTCP(port)
	go udp.RunUDP(port)
	<-make(chan interface{})
}
