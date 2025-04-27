package communication

import (
	"internal"
	"net"
)

type clientConnectionImpl struct {
	connection net.Conn
	delegate   internal.ClientConnectionDelegate
	stopped    bool
}

func createClientConnection(conn net.Conn) clientConnectionImpl {
	return clientConnectionImpl{
		connection: conn,
		delegate:   nil,
		stopped:    false,
	}
}

func (conn *clientConnectionImpl) GetAdressString() string {
	return conn.connection.RemoteAddr().String()
}

func (conn *clientConnectionImpl) SetDelegate(delegate internal.ClientConnectionDelegate) {
	conn.delegate = delegate
}

func (conn *clientConnectionImpl) StartAsync() {
	go conn.startReader()
}

func (conn *clientConnectionImpl) SendMessage(
	id uint64, msgType internal.ServerMessageType, data []byte) {
	// TODO:
}

func (conn *clientConnectionImpl) StopAsync() {
	conn.connection.Close()
}

func (conn *clientConnectionImpl) startReader() {
	// TODO:
}

func ReadMessage() {
	// readingSize := true
	// messageSize := uint64(0)
	// buf := make([]byte, 4096)
	// message := new(bytes.Buffer)
	// for {
	// 	size, err := c.connection.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			log.Printf("Client has been disconnected: %v", c.connection.RemoteAddr())
	// 		} else {
	// 			log.Printf("Client network error: %v", err)
	// 		}
	// 		break
	// 	}
	// 	message.Write(buf[:size])

	// 	// TODO: создать класс соединения, который будет владеть сетевым соединением и разбивать входящие сообщения
	// 	if readingSize {
	// 		if message.Len() < 8 {
	// 			continue
	// 		}
	// 		messageSize = binary.BigEndian.Uint64(message.Bytes())
	// 	}

	// 	if messageSize+8 >= uint64(message.Len()) {

	// 	}
	// }
	// c.connection.Close()
}
