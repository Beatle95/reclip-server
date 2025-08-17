package communication

import (
	"encoding/binary"
	"errors"
	"fmt"
	"internal"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const messageHeaderSize = 24
const minHeaderLen = 16
const writeQueueSize = 100

type networkMessage struct {
	id      uint64
	msgType uint16
	data    []byte
}

// TODO: Add better client error notification.
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
		writeQueue: make(chan networkMessage, writeQueueSize),
	}
}

func createClientConnection(conn net.Conn) clientConnectionImpl {
	return clientConnectionImpl{
		connection: conn,
		writeQueue: make(chan networkMessage, writeQueueSize),
	}
}

func (conn *clientConnectionImpl) GetAdressString() string {
	return conn.connection.RemoteAddr().String()
}

func (conn *clientConnectionImpl) ReadIntroduction() ([]byte, error) {
	conn.connection.SetDeadline(time.Now().Add(time.Second * 15))
	defer conn.connection.SetDeadline(time.Time{})

	lenBuf, err := readNBytes(conn.connection, 8)
	if err != nil {
		return nil, err
	}
	msgLen := binary.BigEndian.Uint64(lenBuf)

	msgBuf, err := readNBytes(conn.connection, msgLen)
	if err != nil {
		return nil, err
	}
	msg, err := parseNetworkMessage(msgBuf)
	if err != nil {
		return nil, err
	}
	if msg.msgType != uint16(internal.ClientIntroduction) {
		return nil, fmt.Errorf("got wrong message type instead of introduction: %d", msg.msgType)
	}

	return msg.data, nil
}

func (conn *clientConnectionImpl) SetUp(
	delegate internal.ClientConnectionDelegate,
	taskRunner internal.EventLoop) {
	conn.delegate = delegate
	conn.taskRunner = taskRunner
}

func (conn *clientConnectionImpl) StartHandlingAsync() {
	conn.stopped.Store(false)
	go conn.readerFunc()
	go conn.writerFunc()
}

func (conn *clientConnectionImpl) DisconnectAndStop() {
	conn.stopped.Store(true)
	conn.connection.Close()
	// We have to add some message in case there is no messages in send queue.
	select {
	case conn.writeQueue <- networkMessage{id: 0, msgType: 0, data: nil}:
	default:
	}
}

func (conn *clientConnectionImpl) SendMessage(
	id uint64, msgType internal.ServerMessageType, data []byte) {
	if conn.stopped.Load() {
		return
	}

	select {
	case conn.writeQueue <- networkMessage{id: id, msgType: uint16(msgType), data: data}:
	default:
		log.Printf("Connection %s write queue is full, it will be disconnected.",
			conn.GetAdressString())
		conn.connection.Close()
	}
}

func (conn *clientConnectionImpl) writerFunc() {
	var prefixBuffer [messageHeaderSize]byte
	for !conn.stopped.Load() {
		msg := <-conn.writeQueue
		if msg.data == nil {
			break
		}
		len := uint64(16 + len(msg.data))
		binary.BigEndian.PutUint64(prefixBuffer[:], len)
		binary.BigEndian.PutUint64(prefixBuffer[8:], msg.id)
		binary.BigEndian.PutUint16(prefixBuffer[16:18], msg.msgType)
		// Warning: bytes from 18 to 24 is reserver for future use.

		_, err1 := conn.connection.Write(prefixBuffer[:])
		if err1 != nil {
			conn.stopped.Store(true)
			break
		}

		_, err2 := conn.connection.Write(msg.data)
		if err2 != nil {
			conn.stopped.Store(true)
			break
		}
	}

	// Trying to prevent some annecessary log messages about full queue upon disconnection.
	popAllMessages(&conn.writeQueue)
}

func (conn *clientConnectionImpl) readerFunc() {
	reassembler := createMessageReassembler()
	buf := make([]byte, 4096)
	for !conn.stopped.Load() {
		size, err := conn.connection.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Client network error: %v", err)
			}
			break
		}

		reassembler.ProcessChunk(buf[:size])
		if reassembler.IsBroken() {
			conn.connection.Close()
			conn.stopped.Store(true)
			break
		}

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

func readNBytes(conn net.Conn, n uint64) ([]byte, error) {
	buffer := make([]byte, n)
	var totalRead uint64 = 0

	for totalRead < n {
		bytesRead, err := conn.Read(buffer[totalRead:])
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}
		if bytesRead < 0 {
			return nil, errors.New("read error")
		}
		totalRead += uint64(bytesRead)
	}

	return buffer, nil
}

func popAllMessages(channel *chan networkMessage) {
	for {
		select {
		case <-*channel:
		default:
			return
		}
	}
}
