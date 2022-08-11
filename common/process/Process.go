package process

type Process = byte

const (
	_ Process = iota
	ServerRegister
	ServerNameExist
	ServerRegisterSuccess

	ClientRegister
	NameNotExist

	ClientAddress

	DoSuccess
	ServerAddress
)

func Message(process Process, src []byte) []byte {
	result := make([]byte, len(src)+1)
	result[0] = process
	copy(result[1:], src)
	return result
}
