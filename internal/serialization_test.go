package internal

import (
	"testing"
)

const expectedSerialized = "{\"ClientId\":1,\"ClientName\":\"name_value\",\"TextData\":[\"text1\",\"text2\"]}"

func TestSerialization(t *testing.T) {
	var dataToSerialize ClientData
	dataToSerialize.Id = 1
	dataToSerialize.Name = "name_value"
	dataToSerialize.Data.Text.PushBack("text1")
	dataToSerialize.Data.Text.PushBack("text2")
	buf := SerializeClientData(&dataToSerialize)
	if buf == nil {
		t.Error("Unable to serialize client data")
	}

	if string(buf) != expectedSerialized {
		t.Errorf("Unexpected serialized value \n1:%s\n2:%s", string(buf), expectedSerialized)
	}

	deserialized, err := DeserializeClientData(buf)
	if err != nil {
		t.Error("Deserialization error")
	}
	if !IsEqual(dataToSerialize, deserialized) {
		t.Error("Deserialized data is not equeal to serialized one")
	}
}
