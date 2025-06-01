package communication

import (
	"bytes"
	"container/list"
	"encoding/binary"
)

const maxMessageLen = 1 * 1024 * 1024 * 1024

type messageReassembler struct {
	messages        *list.List
	buffer          bytes.Buffer
	nextMessageSize uint64
	readingLen      bool
}

func createMessageReassembler() messageReassembler {
	return messageReassembler{
		messages:        list.New(),
		nextMessageSize: 0,
		readingLen:      true,
	}
}

func (a *messageReassembler) ProcessChunk(data []byte) {
	// TODO: Protect from buffer overflow.
	a.buffer.Write(data)
	continue_process := true
	for continue_process {
		if a.readingLen {
			continue_process = a.tryReadLen()
		} else {
			continue_process = a.tryReadMessage()
		}
	}
}

func (a *messageReassembler) HasMessage() bool {
	return a.messages.Len() > 0
}

func (a *messageReassembler) PopMessage() []byte {
	return a.messages.Remove(a.messages.Front()).([]byte)
}

// Returns true if length was read.
func (a *messageReassembler) tryReadLen() bool {
	if a.buffer.Len() < 8 {
		return false
	}
	lenBuf := make([]byte, 8)
	bytesRead, err := a.buffer.Read(lenBuf)
	if len(lenBuf) != bytesRead || err != nil {
		panic("Something went wrong while processing incoming message length")
	}
	a.nextMessageSize = binary.BigEndian.Uint64(lenBuf)
	if a.nextMessageSize > maxMessageLen {
		// TODO: Enter broken state.
		panic("We have to process this")
	}
	a.readingLen = false
	return true
}

// Returns true if we have read the message.
func (a *messageReassembler) tryReadMessage() bool {
	if uint64(a.buffer.Len()) < a.nextMessageSize {
		return false
	}

	msgData := make([]byte, a.nextMessageSize)
	bytesRead, err := a.buffer.Read(msgData)
	if bytesRead != len(msgData) || err != nil {
		panic("Something went wrong while processing incoming message data")
	}
	a.nextMessageSize = 0
	a.readingLen = true
	a.messages.PushBack(msgData)
	return true
}
