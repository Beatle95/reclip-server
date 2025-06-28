package communication

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type TestData struct {
	data_with_len []byte
	data_only     []byte
}

func TestCorrectMessages(t *testing.T) {
	test_data := createTestBuffer(2048)
	assembler := createMessageReassembler()

	assembler.ProcessChunk(test_data.data_with_len)
	if !assembler.HasMessage() {
		t.FailNow()
	}
	msg := assembler.PopMessage()
	if msg == nil || !bytes.Equal(msg, test_data.data_only) || assembler.HasMessage() {
		t.FailNow()
	}

	for step_size := 1; step_size < 128; step_size += 8 {
		var offset int
		for offset = 0; offset+step_size < len(test_data.data_with_len); offset += step_size {
			assembler.ProcessChunk(test_data.data_with_len[offset : offset+step_size])
		}
		assembler.ProcessChunk(test_data.data_with_len[offset:])
		if !assembler.HasMessage() {
			t.Fatalf("Expected to have message. Step size was: %d", step_size)
		}
		msg := assembler.PopMessage()
		if msg == nil || !bytes.Equal(msg, test_data.data_only) || assembler.HasMessage() {
			t.Fatalf("Step size was: %d", step_size)
		}
	}

	const messages_count = 10
	for i := 0; i < messages_count; i++ {
		for offset := 0; offset < len(test_data.data_with_len); offset++ {
			assembler.ProcessChunk(test_data.data_with_len[offset : offset+1])
		}
	}
	for i := 0; i < messages_count; i++ {
		if !assembler.HasMessage() {
			t.Fatal("Expected to have message")
		}
		msg := assembler.PopMessage()
		if msg == nil || !bytes.Equal(msg, test_data.data_only) {
			t.Fatal("Incorrect result message")
		}
	}
	if assembler.HasMessage() {
		t.Fatal("Unexpected message at the end")
	}
}

// TODO: Test zero len message. Test too long message.

func createTestBuffer(desired_data_len uint64) TestData {
	data := make([]byte, desired_data_len+8)
	binary.BigEndian.PutUint64(data, desired_data_len)

	result := TestData{
		data_with_len: data,
		data_only:     data[8:],
	}

	var ch byte = 0
	for i := 0; i < len(result.data_only); i++ {
		result.data_only[i] = ch
		ch += 1
	}
	return result
}
