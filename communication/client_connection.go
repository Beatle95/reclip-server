package communication

import (
	"encoding/binary"
	"internal"
	"io"
	"log"
	"net"
	"sync/atomic"
)

const messageHeaderSize = 18

type networkMessage struct {
	id      uint64
	msgType uint16
	data    []byte
}

type clientConnectionImpl struct {
	connection net.Conn
	delegate   internal.ClientConnectionDelegate
	stopped    atomic.Bool

	writeQueue chan networkMessage
}

func createClientConnection(conn net.Conn) clientConnectionImpl {
	return clientConnectionImpl{
		connection: conn,
		delegate:   nil,
		writeQueue: make(chan networkMessage, 20),
	}
}

func (conn *clientConnectionImpl) GetAdressString() string {
	return conn.connection.RemoteAddr().String()
}

func (conn *clientConnectionImpl) SetDelegate(delegate internal.ClientConnectionDelegate) {
	conn.delegate = delegate
}

func (conn *clientConnectionImpl) StartAsync() {
	go conn.readerFunc()
	go conn.writerFunc()
}

func (conn *clientConnectionImpl) StopAsync() {
	conn.stopped.Store(true)
	conn.connection.Close()
	// We have to add some message in case there is no messages in send queue.
	conn.writeQueue <- networkMessage{id: 0, msgType: 0, data: nil}
	// TODO: join goroutines.
}

func (conn *clientConnectionImpl) SendMessage(
	id uint64, msgType internal.ServerMessageType, data []byte) {
	conn.writeQueue <- networkMessage{id: id, msgType: uint16(msgType), data: data}
}

func (conn *clientConnectionImpl) writerFunc() {
	var prefixBuffer [messageHeaderSize]byte
	for !conn.stopped.Load() {
		msg := <-conn.writeQueue
		len := uint64(8 + 2 + len(msg.data))
		binary.LittleEndian.PutUint64(prefixBuffer[:], len)
		binary.LittleEndian.PutUint64(prefixBuffer[8:], msg.id)
		binary.LittleEndian.PutUint16(prefixBuffer[16:], msg.msgType)

		_, err1 := conn.connection.Write(prefixBuffer[:])
		if err1 != nil {
			conn.stopped.Store(true)
			return
		}

		_, err2 := conn.connection.Write(msg.data)
		if err2 != nil {
			conn.stopped.Store(true)
			return
		}
	}
}

func (conn *clientConnectionImpl) readerFunc() {
	var reassembler messageReassembler
	buf := make([]byte, 4096)
	for !conn.stopped.Load() {
		size, err := conn.connection.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client has been disconnected: %v", conn.connection.RemoteAddr())
			} else {
				log.Printf("Client network error: %v", err)
			}
			break
		}

		reassembler.ProcessChunk(buf[:size])
		for reassembler.HasMessage() {
			msg := reassembler.PopMessage()
			conn.delegate.ProcessMessage(msg.id, internal.ClientMessageType(msg.msgType), msg.data)
		}
	}
	conn.delegate.OnDisconnected()
}
