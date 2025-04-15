package internal

import "fmt"

type Message struct {
	id      uint64
	msgType MessageType
	data    []byte
}

type ClientConnection interface {
	SendMessage()
	ReadMessage() (Message, error)
	Disconnect()
}

type ClientDelegate interface {
	OnClientDisconnected(id string)
}

type Client interface {
	GetPublicId() string
	GetName() string
	HandleConnection()

	// TODO: Remote client functions.
	// NotifyClientConnected(id string)
	// NotifyClientDisconnected(id string)
	// NotifyTextAdded(id string, text string)
	// NotifyClientSynced()
}

type clientImpl struct {
	connection ClientConnection
	delegate   ClientDelegate
	data       ClientData
}

func CreateClient(
	connection ClientConnection,
	delegate ClientDelegate,
	publicId string,
	name string) Client {
	return &clientImpl{
		connection: connection,
		delegate:   delegate,
		data: ClientData{
			id:   publicId,
			name: name,
		},
	}
}

func (c *clientImpl) GetPublicId() string {
	return c.data.id
}

func (c *clientImpl) GetName() string {
	return c.data.name
}

func (c *clientImpl) HandleConnection() {
	for {
		_, err := c.connection.ReadMessage()
		if err != nil {
			fmt.Printf("Client will be disconnected: %v", err)
			break
		}
		// TODO: deal with msg
	}

	c.connection.Disconnect()
	c.delegate.OnClientDisconnected(c.data.id)
}
