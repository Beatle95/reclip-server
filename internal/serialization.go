package internal

import (
	"encoding/json"
	"fmt"
)

type serverIntroductionJson struct {
	Version string
}

type syncJson struct {
	ThisHostData clientJson
	OtherData    []clientJson
}

type clientJson struct {
	ClientId   uint64
	ClientName string
	TextData   []string
}

type clientIdJson struct {
	ClientId uint64
}

type textUpdateJson struct {
	ClientId uint64
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
	otherDataJson := make([]clientJson, len(otherData))
	for index, elem := range otherData {
		otherDataJson[index] = clientDataToJsonData(&elem)
	}

	data, err := json.Marshal(syncJson{
		ThisHostData: clientDataToJsonData(&thisData),
		OtherData:    otherDataJson,
	})
	if err != nil {
		return nil
	}
	return data
}

func SerializeClientData(clientData *ClientData) []byte {
	data, err := json.Marshal(clientDataToJsonData(clientData))
	if err != nil {
		return nil
	}
	return data
}

func SerializeClientId(id uint64) []byte {
	data, err := json.Marshal(clientIdJson{ClientId: id})
	if err != nil {
		return nil
	}
	return data
}

func SerializeTextUpdate(id uint64, text string) []byte {
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

func DeserializeClientId(data []byte) (uint64, error) {
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

func clientDataToJsonData(clientData *ClientData) clientJson {
	client := clientJson{
		ClientId:   clientData.Id,
		ClientName: clientData.Name,
		TextData:   make([]string, clientData.Data.Text.Len()),
	}
	for i := 0; i < clientData.Data.Text.Len(); i++ {
		client.TextData[i] = clientData.Data.Text.At(i)
	}
	return client
}
