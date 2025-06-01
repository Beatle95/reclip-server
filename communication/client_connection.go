package communication

import (
	"encoding/binary"
	"errors"
	"internal"
	"io"
	"log"
	"net"
	"sync/atomic"
)

const messageHeaderSize = 24
const minHeaderLen = 16

type networkMessage struct {
	id      uint64
	msgType uint16
	data    []byte
}

type clientConnectionImpl struct {
	connection net.Conn
	delegate   internal.ClientConnectionDelegate
	taskRunner internal.EventLoop
	stopped    atomic.Bool

	writeQueue chan networkMessage
}

func CreateClientConnectionForTesting(conn net.Conn) internal.ClientConnection {
	return &clientConnectionImpl{
		connection: conn,
		delegate:   nil,
		writeQueue: make(chan networkMessage, 20),
	}
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

func (conn *clientConnectionImpl) SetUp(
	delegate internal.ClientConnectionDelegate,
	taskRunner internal.EventLoop) {
	conn.delegate = delegate
	conn.taskRunner = taskRunner
}

func (conn *clientConnectionImpl) StartAsync() {
	conn.stopped.Store(false)
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
		len := uint64(16 + len(msg.data))
		binary.BigEndian.PutUint64(prefixBuffer[:], len)
		binary.BigEndian.PutUint64(prefixBuffer[8:], msg.id)
		binary.BigEndian.PutUint16(prefixBuffer[16:18], msg.msgType)
		// Warning: bytes from 18 to 24 is reserver for future use.

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
	reassembler := createMessageReassembler()
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
			msg, err := parseNetworkMessage(reassembler.PopMessage())
			if err == nil {
				conn.taskRunner.PostTask(func() {
					conn.delegate.ProcessMessage(msg.id, internal.ClientMessageType(msg.msgType), msg.data)
				})
			} else {
				log.Printf("Error parsing network header: %s", err.Error())
			}
		}
	}
	conn.taskRunner.PostTask(conn.delegate.OnDisconnected)
}

func parseNetworkMessage(data []byte) (networkMessage, error) {
	var result networkMessage
	if len(data) < minHeaderLen {
		return result, errors.New("some messages was skipped because it was too short")
	}
	result.id = binary.BigEndian.Uint64(data[:8])
	result.msgType = binary.BigEndian.Uint16(data[8:10])
	// Warning: bytes from 10 to 16 is reserver for future use.
	result.data = data[16:]
	return result, nil
}
