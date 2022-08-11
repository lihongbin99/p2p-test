package main

import (
	"flag"
	"p2p-test/client/tcp"
	"p2p-test/client/udp"
	"p2p-test/common/logger"
	"time"
)

var (
	log  = logger.Log("Client Main")
	name string

	cType string

	test       string
	testModel  string
	serverAddr string
)

func init() {
	flag.StringVar(&name, "name", "12345678", "p2p name, minLen:6, maxLen:20")

	flag.StringVar(&cType, "type", "c", "client type: [server, s, client, c]")

	flag.StringVar(&test, "test", "u", "test model: [tcp, t, udp, u]")
	flag.StringVar(&testModel, "model", "f", "test model: [full-cone, f, restricted-cone, r, symmetric, s]")
	//flag.StringVar(&serverAddr, "server", "0.0.0.0:13520", "server addr: [ip:port]")
	flag.StringVar(&serverAddr, "server", "18.141.239.196:13520", "server addr: [ip:port]")
	flag.Parse()
}

func main2() {
	exit := make(chan interface{})

	go func() {
		udp.RunServerUDP(name, serverAddr, testModel)
		exit <- 1
	}()

	time.Sleep(1 * time.Second)

	go func() {
		udp.RunClientUDP(name, serverAddr, testModel)
		exit <- 1
	}()

	_ = <-exit
	_ = <-exit
}

func main() {
	if cType == "server" || cType == "s" {
		if test == "tcp" || test == "t" {
			tcp.RunServerTCP(name, serverAddr, testModel)
		} else if test == "udp" || test == "u" {
			udp.RunServerUDP(name, serverAddr, testModel)
		} else {
			log.Fatal("test model error(", test, ")not in [tcp, t, udp, u]")
		}
	} else if cType == "client" || cType == "c" {
		if test == "tcp" || test == "t" {
			tcp.RunClientTCP(name, serverAddr, testModel)
		} else if test == "udp" || test == "u" {
			udp.RunClientUDP(name, serverAddr, testModel)
		} else {
			log.Fatal("test model error(", test, ")not in [tcp, t, udp, u]")
		}
	} else {
		log.Fatal("client type error(", cType, ")not in [server, s, client, c]")
	}
}
