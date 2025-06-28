package communication

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
)

const maxMessageLen = 1 * 1024 * 1024 * 1024

type messageReassembler struct {
	messages        *list.List
	buffer          bytes.Buffer
	nextMessageSize uint64
	readingLen      bool
	isBroken        bool
}

func createMessageReassembler() messageReassembler {
	return messageReassembler{
		messages:        list.New(),
		nextMessageSize: 0,
		readingLen:      true,
		isBroken:        false,
	}
}

func (a *messageReassembler) ProcessChunk(data []byte) {
	if a.isBroken {
		panic("Trying to process chunk after reassembler is broken")
	}

	a.buffer.Write(data)
	continueProcess := true
	for continueProcess {
		if a.readingLen {
			continueProcess = a.tryReadLen()
		} else {
			continueProcess = a.tryReadMessage()
		}
	}

	if a.buffer.Len() > maxMessageLen {
		// Buffer leftover length is constrained by message lendgth, which is constrained by
		// maxMessageLen, so this situation must not happen.
		a.isBroken = true
		fmt.Printf("buffer leftover length was too long '%d'", a.buffer.Len())
	}
}

func (a *messageReassembler) HasMessage() bool {
	return a.messages.Len() > 0
}

func (a *messageReassembler) PopMessage() []byte {
	return a.messages.Remove(a.messages.Front()).([]byte)
}

func (a *messageReassembler) IsBroken() bool {
	return a.isBroken
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
		a.isBroken = true
		return false
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
