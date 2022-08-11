package udp

import (
	"fmt"
	"net"
	"p2p-test/common/logger"
	"p2p-test/common/process"
)

var (
	log       = logger.Log("UDP Server")
	serverMap = make(map[string]*NameUDP)
	clientMap = make(map[string]*NameUDP)
)

type NameUDP struct {
	*net.UDPAddr
	wait chan interface{}
}

func RunUDP(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("resolve server addr", err)
	}
	udp, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("listen server", err)
	}
	log.Info("listen server", udp.LocalAddr().String())

	buf := make([]byte, 32)
	for {
		if readLen, remoteAddr, err := udp.ReadFromUDP(buf); err != nil {
			log.Error("read from udp", err)
		} else {
			cp := make([]byte, readLen)
			copy(cp, buf[:readLen])
			go doRunUDP(udp, cp, remoteAddr)
		}
	}
}

func doRunUDP(udp *net.UDPConn, buf []byte, addr *net.UDPAddr) {
	name := string(buf[1:])
	switch buf[0] {
	case process.ServerRegister:
		if _, exist := serverMap[name]; exist {
			_, _ = udp.WriteToUDP(process.Message(process.ServerNameExist, []byte(name)), addr)
			return
		}
		tcpName := NameUDP{UDPAddr: addr, wait: make(chan interface{})}
		serverMap[name] = &tcpName
		go func(tcpName *NameUDP) {
			_ = <-tcpName.wait
			defer delete(serverMap, name)
			log.Info("server exit", addr.String(), "name", name)
		}(&tcpName)
		log.Info("server register success", addr.String(), "name", name)

		_, _ = udp.WriteToUDP(process.Message(process.ServerRegisterSuccess, []byte(name)), addr)

	case process.ClientRegister:
		serverUDP, exist := serverMap[name]
		if !exist {
			_, _ = udp.WriteToUDP(process.Message(process.NameNotExist, []byte(name)), addr)
			return
		}
		tcpName := NameUDP{UDPAddr: addr, wait: make(chan interface{})}
		clientMap[name] = &tcpName
		go func(tcpName *NameUDP) {
			_ = <-tcpName.wait
			defer delete(clientMap, name)
			log.Info("client exit", addr.String(), "name", name)
		}(&tcpName)
		log.Info("client register success", addr.String(), "name", name)

		_, _ = udp.WriteToUDP(process.Message(process.ClientAddress, []byte(addr.String())), serverUDP.UDPAddr)
		log.Info("client -> server", addr.String(), " -> ", serverUDP.String())

	case process.DoSuccess:
		if clientUDP, exist := clientMap[name]; exist {
			_, _ = udp.WriteToUDP(process.Message(process.ServerAddress, []byte(addr.String())), clientUDP.UDPAddr)
			log.Info("server -> client", addr.String(), " -> ", clientUDP.String())
			clientUDP.wait <- 1
		}
		if serverUDP, exist := serverMap[name]; exist {
			serverUDP.wait <- 1
		}
	}
}
