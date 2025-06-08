package internal

type ClientConnectionDelegate interface {
	OnDisconnected()
	ProcessMessage(id uint64, msgType ClientMessageType, data []byte)
}

type ClientConnection interface {
	GetAdressString() string
	ReadIntroduction() ([]byte, error)
	SetUp(delegate ClientConnectionDelegate, taskRunner EventLoop)

	StartHandlingAsync()
	DisconnectAndStop()

	SendMessage(id uint64, msgType ServerMessageType, data []byte)
}
