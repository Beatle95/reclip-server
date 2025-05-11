package internal

type ClientConnectionDelegate interface {
	OnDisconnected()
	ProcessMessage(id uint64, msgType ClientMessageType, data []byte)
}

type ClientConnection interface {
	GetAdressString() string
	SetUp(delegate ClientConnectionDelegate, taskRunner EventLoop)

	StartAsync()
	StopAsync()

	SendMessage(id uint64, msgType ServerMessageType, data []byte)
}
