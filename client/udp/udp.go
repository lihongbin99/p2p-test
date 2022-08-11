package udp

import "p2p-test/common/logger"

var (
	log = logger.Log("Client-UDP")
	buf = make([]byte, 64)
)
