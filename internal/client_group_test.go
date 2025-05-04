package internal

import (
	"testing"
)

type OthersTextData struct {
	id   string
	text string
}

type MockClient struct {
	data                     ClientData
	othersText               []OthersTextData
	handleConnected          uint32
	notifyClientConnected    uint32
	notifyClientDisconnected uint32
	notifyClientSynced       uint32
}

func (c *MockClient) GetPublicId() string {
	return c.data.Id
}

func (c *MockClient) GetName() string {
	return c.data.Name
}

func (c *MockClient) GetClientData() *ClientData {
	return &c.data
}

func (c *MockClient) HandleConnection() {
	c.handleConnected++
}

func (c *MockClient) NotifyClientConnected(id string) {
	c.notifyClientConnected++
}

func (c *MockClient) NotifyClientDisconnected(id string) {
	c.notifyClientDisconnected++
}

func (c *MockClient) NotifyTextAdded(id string, text string) {
	if c.othersText == nil {
		c.othersText = make([]OthersTextData, 0)
	}
	c.othersText = append(c.othersText, OthersTextData{id: id, text: text})
}

func (c *MockClient) NotifyClientSynced(data *ClientData) {
	c.notifyClientSynced++
}

func TestClientGroup(t *testing.T) {
	client1 := MockClient{
		data: ClientData{
			Id:   "id1",
			Name: "name1",
		},
	}
	client2 := MockClient{
		data: ClientData{
			Id:   "id2",
			Name: "name2",
		},
	}
	client3 := MockClient{
		data: ClientData{
			Id:   "id3",
			Name: "name3",
		},
	}

	testGroup := CreateClientGroup()
	testGroup.HandleNewClient(&client1)
	if client1.handleConnected != 1 {
		t.Errorf("Handle connected not called for client 1")
	}

	testGroup.HandleNewClient(&client2)
	if client2.handleConnected != 1 {
		t.Errorf("Handle connected not called for client 2")
	}
	if client1.notifyClientConnected != 1 {
		t.Errorf("Notify connected is not called after 2 connected")
	}

	testGroup.HandleNewClient(&client3)
	if client3.handleConnected != 1 {
		t.Errorf("Handle connected not called for client 3")
	}
	if client1.notifyClientConnected != 2 || client2.notifyClientConnected != 1 {
		t.Errorf("Notify connected is not called after 3 connected")
	}

	testGroup.OnClientSynced(&client1)
	if client1.notifyClientSynced != 0 || client2.notifyClientSynced != 1 ||
		client3.notifyClientSynced != 1 {
		t.Error("OnClientSynced processed incorrectly")
	}

	sometText := "some added text"
	client1.data.Data.Text.PushBack(sometText)
	testGroup.OnTextAdded(&client1, sometText)
	if len(client2.othersText) != 1 || client2.othersText[0].id != client1.GetPublicId() ||
		client2.othersText[0].text != sometText {
		t.Error("Wrong text added processing")
	}

	syncData := testGroup.GetFullSyncData(&client3)
	if !IsEqual(syncData[0], client1.data) || !IsEqual(syncData[1], client2.data) {
		t.Errorf("Bad full sync response, sync data was: %v", syncData)
	}

	testGroup.OnClientDisconnected(&client3)
	if client1.notifyClientDisconnected != 1 || client2.notifyClientDisconnected != 1 {
		t.Error("OnClientDisconnected processed incorrectly")
	}
}
