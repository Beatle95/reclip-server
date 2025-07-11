package internal

import (
	"testing"
)

type OthersTextData struct {
	id   uint64
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

type MockClientConnection struct{}

func (c *MockClient) IsConnected() bool {
	return c.handleConnected != 0
}

func (c *MockClient) GetClientData() *ClientData {
	return &c.data
}

func (c *MockClient) HandleConnection(connection ClientConnection) {
	c.handleConnected++
}

func (c *MockClient) NotifyClientConnected(id uint64) {
	c.notifyClientConnected++
}

func (c *MockClient) NotifyClientDisconnected(id uint64) {
	c.notifyClientDisconnected++
}

func (c *MockClient) NotifyTextAdded(id uint64, text string) {
	if c.othersText == nil {
		c.othersText = make([]OthersTextData, 0)
	}
	c.othersText = append(c.othersText, OthersTextData{id: id, text: text})
}

func (c *MockClient) NotifyClientSynced(data *ClientData) {
	c.notifyClientSynced++
}

func (c *MockClientConnection) GetAdressString() string { return "" }

func (c *MockClientConnection) ReadIntroduction() ([]byte, error) { return nil, nil }

func (c *MockClientConnection) SetUp(delegate ClientConnectionDelegate, taskRunner EventLoop) {}

func (c *MockClientConnection) StartHandlingAsync() {}

func (c *MockClientConnection) DisconnectAndStop() {}

func (c *MockClientConnection) SendMessage(id uint64, msgType ServerMessageType, data []byte) {}

func TestClientGroup(t *testing.T) {
	client1 := MockClient{
		data: ClientData{
			Id:   1,
			Name: "name1",
		},
	}
	client2 := MockClient{
		data: ClientData{
			Id:   2,
			Name: "name2",
		},
	}
	client3 := MockClient{
		data: ClientData{
			Id:   3,
			Name: "name3",
		},
	}

	testGroup := CreateClientGroup()
	testGroup.AddClient(&client1)
	testGroup.AddClient(&client2)
	testGroup.AddClient(&client3)

	testGroup.HandleConnection(1, &MockClientConnection{})
	testGroup.GetTaskRunner().RunUntilIdle()
	if client1.handleConnected != 1 {
		t.Errorf("Handle connected not called for client 1")
	}
	if client1.notifyClientConnected != 0 || client2.notifyClientConnected != 1 ||
		client3.notifyClientConnected != 1 {
		t.Errorf("Notify connected is not called after 1 connected")
	}

	testGroup.HandleConnection(2, &MockClientConnection{})
	testGroup.GetTaskRunner().RunUntilIdle()
	if client2.handleConnected != 1 {
		t.Errorf("Handle connected not called for client 2")
	}
	if client1.notifyClientConnected != 1 || client2.notifyClientConnected != 1 ||
		client3.notifyClientConnected != 2 {
		t.Errorf("Notify connected is not called after 2 connected")
	}

	testGroup.HandleConnection(3, &MockClientConnection{})
	testGroup.GetTaskRunner().RunUntilIdle()
	if client3.handleConnected != 1 {
		t.Errorf("Handle connected not called for client 3")
	}
	if client1.notifyClientConnected != 2 || client2.notifyClientConnected != 2 ||
		client3.notifyClientConnected != 2 {
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
	if len(client2.othersText) != 1 || client2.othersText[0].id != client1.GetClientData().Id ||
		client2.othersText[0].text != sometText {
		t.Error("Wrong text added processing")
	}

	syncData := testGroup.GetFullSyncData(&client3)
	if !IsEqual(syncData[0], client1.data) {
		t.Errorf("Bad full sync response, sync[0] data was: %v", syncData[0])
	}
	if !IsEqual(syncData[1], client2.data) {
		t.Errorf("Bad full sync response, sync[1] data was: %v", syncData[1])
	}

	testGroup.OnClientDisconnected(&client3)
	if client1.notifyClientDisconnected != 1 || client2.notifyClientDisconnected != 1 {
		t.Error("OnClientDisconnected processed incorrectly")
	}
}
