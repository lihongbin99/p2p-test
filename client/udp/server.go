package udp

import (
	"net"
	"p2p-test/common/process"
)

func RunServerUDP(name string, serverAddrS string, model string) {
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

	if _, err = serverUDP.Write(process.Message(process.ServerRegister, []byte(name))); err != nil {
		log.Fatal("write server register error", err)
	}

	readLen, err := serverUDP.Read(buf)
	if err != nil {
		log.Fatal("read register status error", err)
	}
	switch buf[0] {
	case process.ServerNameExist:
		log.Fatal("server name exist", string(buf[1:readLen]))
	case process.ServerRegisterSuccess:
		log.Info("server register success", string(buf[1:readLen]))
	}

	readLen, err = serverUDP.Read(buf)
	if err != nil {
		log.Fatal("read remote addr error", err)
	}
	switch buf[0] {
	case process.ClientAddress:
	default:
		log.Fatal("read message no remote address", buf)
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", string(buf[1:readLen]))
	if err != nil {
		log.Fatal("resolve remote addr error", err)
	}
	log.Info("remote address", remoteAddr.String())

	// TODO model
	_ = serverUDP.Close()
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		log.Fatal("dial remote error", err)
	}
	if _, err = conn.Write([]byte("test message")); err != nil {
		log.Fatal("write test message error", err)
	}
	_ = conn.Close()
	serverUDP, _ = net.DialUDP("udp", localAddr, serverAddr)

	if _, err = serverUDP.Write(process.Message(process.DoSuccess, []byte(name))); err != nil {
		log.Fatal("notify client do success error", err)
	}
	_ = serverUDP.Close()

	udp, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal("listen udp error", err)
	}

	readLen, remoteAddr, err = udp.ReadFromUDP(buf)
	if err != nil {
		log.Fatal("read remote message error", err)
	}
	log.Info("addr", remoteAddr.String(), "message", string(buf[:readLen]))
	_, _ = udp.WriteToUDP([]byte("Hello I am Server"), remoteAddr)
}
