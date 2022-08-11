package udp

import (
	"net"
	"p2p-test/common/process"
	"time"
)

func RunClientUDP(name string, serverAddrS string, model string) {
	serverAddr, err := net.ResolveUDPAddr("udp", serverAddrS)
	if err != nil {
		log.Fatal("resolve server addr error", err)
	}
	serverUDP, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal("dial server error", err)
	}
	localAddrS := serverUDP.LocalAddr().String()
	localAddr, _ := net.ResolveUDPAddr("udp", localAddrS)
	log.Info("local address", localAddr.String())

	if _, err = serverUDP.Write(process.Message(process.ClientRegister, []byte(name))); err != nil {
		log.Fatal("write client register error", err)
	}

	readLen, err := serverUDP.Read(buf)
	if err != nil {
		log.Fatal("read remote addr error", err)
	} else if readLen == 1 {
		log.Fatal("read remote addr length = 1")
	}
	switch buf[0] {
	case process.NameNotExist:
		log.Fatal("name not exist", string(buf[1:readLen]))
	case process.ServerAddress:
	default:
		log.Fatal("read message no remote address", buf)
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", string(buf[1:readLen]))
	if err != nil {
		log.Fatal("resolve remote addr error", err)
	}
	log.Info("remote address", remoteAddr.String())

	// TODO model

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		log.Fatal("dial remote error", err)
	}

	if _, err = conn.Write([]byte("Hello I am Client")); err != nil {
		log.Fatal("write to remote error", err)
	}
	log.Info("write to remote success")

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	readLen, err = conn.Read(buf)
	if err != nil {
		// TODO 可能用防火墙, 所以再发一次
		_ = conn.SetReadDeadline(time.Time{})

		if _, err = conn.Write([]byte("Hello I am Client")); err != nil {
			log.Fatal("write to remote error", err)
		}
		log.Info("write to remote success")
		readLen, err = conn.Read(buf)
		if err != nil {
			log.Fatal("read remote message error", err)
		}
	}
	log.Info("message", string(buf[:readLen]))
}
