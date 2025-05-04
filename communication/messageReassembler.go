package communication

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"log"
)

type messageReassembler struct {
	messages        list.List
	buffer          bytes.Buffer
	nextMessageSize uint64
	readingLen      bool
}

func (a *messageReassembler) ProcessChunk(data []byte) {
	a.buffer.Write(data)
	a.process()
}

func (a *messageReassembler) HasMessage() bool {
	return a.messages.Len() > 0
}

func (a *messageReassembler) PopMessage() networkMessage {
	firstElem := a.messages.Front()
	return a.messages.Remove(firstElem).(networkMessage)
}

func (a *messageReassembler) process() {
	if a.readingLen {
		a.tryReadLen()
	} else {
		a.tryReadMessage()
	}
}

func (a *messageReassembler) tryReadLen() {
	if a.buffer.Len() < 8 {
		return
	}
	lenBuf := make([]byte, 8)
	bytesRead, err := a.buffer.Read(lenBuf)
	if len(lenBuf) != bytesRead || err != nil {
		panic("Something wen wrong while processing incoming message length")
	}
	a.nextMessageSize = binary.BigEndian.Uint64(lenBuf)
}

func (a *messageReassembler) tryReadMessage() {
	if a.nextMessageSize < uint64(a.buffer.Len()) {
		return
	}

	msgData := make([]byte, a.nextMessageSize)
	bytesRead, err := a.buffer.Read(msgData)
	if bytesRead != len(msgData) || err != nil {
		panic("Something went wrong while processing incoming message data")
	}
	if bytesRead < 10 {
		log.Printf("Some messages was skipped because it was too short")
		return
	}

	a.messages.PushBack(networkMessage{
		id:      binary.BigEndian.Uint64(msgData[:]),
		msgType: binary.BigEndian.Uint16(msgData[8:]),
		data:    msgData[10:],
	})
}
