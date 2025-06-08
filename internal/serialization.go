package internal

import (
	"encoding/json"
	"fmt"
)

type clientIntroductionJson struct {
	Secret string
}

type serverIntroductionJson struct {
	Version string
}

type syncJson struct {
	ThisHostData ClientData
	OtherData    []ClientData
}

type clientJson struct {
	ClientId   string
	ClientName string
	TextData   []string
}

type clientIdJson struct {
	ClientId string
}

type textUpdateJson struct {
	ClientId string
	Text     string
}

type textJson struct {
	Text string
}

type errorJson struct {
	ErrorText string
}

func SerializeIntroduction(ver Version) []byte {
	data, err := json.Marshal(serverIntroductionJson{
		Version: fmt.Sprintf("%d.%d.%d", ver.major, ver.minor, ver.build_number),
	})
	if err != nil {
		return nil
	}
	return data
}

func SerializeSync(thisData ClientData, otherData []ClientData) []byte {
	data, err := json.Marshal(syncJson{
		ThisHostData: thisData,
		OtherData:    otherData,
	})
	if err != nil {
		return nil
	}
	return data
}

func SerializeClientData(clientData *ClientData) []byte {
	client := clientJson{
		ClientId:   clientData.Id,
		ClientName: clientData.Name,
		TextData:   make([]string, clientData.Data.Text.Len()),
	}
	for i := 0; i < clientData.Data.Text.Len(); i++ {
		client.TextData[i] = clientData.Data.Text.At(i)
	}

	data, err := json.Marshal(client)
	if err != nil {
		return nil
	}
	return data
}

func SerializeClientId(id string) []byte {
	data, err := json.Marshal(clientIdJson{ClientId: id})
	if err != nil {
		return nil
	}
	return data
}

func SerializeTextUpdate(id string, text string) []byte {
	data, err := json.Marshal(textUpdateJson{ClientId: id, Text: text})
	if err != nil {
		return nil
	}
	return data
}

func SerializeError(errorText string) []byte {
	data, err := json.Marshal(errorJson{ErrorText: errorText})
	if err != nil {
		return nil
	}
	return data
}

func DeserializeClientIntroduction(data []byte) (string, error) {
	var intro clientIntroductionJson
	err := json.Unmarshal(data, &intro)
	return intro.Secret, err
}

func DeserializeClientId(data []byte) (string, error) {
	var clientId clientIdJson
	err := json.Unmarshal(data, &clientId)
	return clientId.ClientId, err
}

func DeserializeText(data []byte) (string, error) {
	var text textJson
	err := json.Unmarshal(data, &text)
	return text.Text, err
}

func DeserializeClientData(data []byte) (ClientData, error) {
	var client clientJson
	err := json.Unmarshal(data, &client)
	clientData := ClientData{
		Id:   client.ClientId,
		Name: client.ClientName,
	}
	for _, val := range client.TextData {
		clientData.Data.Text.PushBack(val)
	}
	return clientData, err
}
