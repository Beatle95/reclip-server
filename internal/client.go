package internal

import "fmt"

type ClientDelegate interface {
	GetTaskRunner() EventLoop
	OnClientDisconnected(client Client)
	GetFullSyncData(syncExcluded Client) []ClientData
	GetClientSyncData(id uint64) *ClientData
	OnTextAdded(client Client, text string)
	OnClientSynced(client Client)
}

type Client interface {
	IsConnected() bool
	GetClientData() *ClientData
	HandleConnection(connection ClientConnection)

	NotifyClientConnected(id uint64)
	NotifyClientDisconnected(id uint64)
	NotifyTextAdded(id uint64, text string)
	NotifyClientSynced(data *ClientData)
}

type clientImpl struct {
	connection ClientConnection
	delegate   ClientDelegate
	data       ClientData
	idCounter  uint64
}

func CreateClient(
	delegate ClientDelegate,
	publicId uint64,
	name string) Client {
	return &clientImpl{
		delegate: delegate,
		data: ClientData{
			Id:   publicId,
			Name: name,
		},
		idCounter: 0,
	}
}

// ClientConnectionDelegate implementations:

func (c *clientImpl) ProcessMessage(id uint64, msgType ClientMessageType, data []byte) {
	switch msgType {
	case ClientResponse:
		// Not used right now.
	case FullSyncRequest:
		c.processFullSyncRequest(id)
	case HostSyncRequest:
		c.processHostSyncRequest(id, data)
	case HostTextUpdate:
		c.processHostTextUpdate(data)
	case SyncThisHost:
		c.processSyncClient(data)
	}
}

func (c *clientImpl) OnDisconnected() {
	c.delegate.OnClientDisconnected(c)
	c.connection = nil
}

// Client implementations:

func (c *clientImpl) IsConnected() bool {
	return c.connection != nil
}

func (c *clientImpl) GetClientData() *ClientData {
	return &c.data
}

func (c *clientImpl) HandleConnection(connection ClientConnection) {
	if c.connection != nil {
		panic("Resetting connection which was already set")
	}
	c.connection = connection
	c.connection.SetUp(c, c.delegate.GetTaskRunner())

	// Right away schedule introduction sending.
	serialized := SerializeIntroduction(GetApplicationVersion())
	c.connection.SendMessage(c.idCounter, ServerIntroduction, serialized)
	c.idCounter++

	// Start connection handling.
	c.connection.StartHandlingAsync()
}

func (c *clientImpl) NotifyClientConnected(id uint64) {
	if c.connection == nil {
		return
	}
	serialized := SerializeClientId(id)
	c.connection.SendMessage(c.idCounter, HostConnected, serialized)
	c.idCounter++
}

func (c *clientImpl) NotifyClientDisconnected(id uint64) {
	if c.connection == nil {
		return
	}
	serialized := SerializeClientId(id)
	c.connection.SendMessage(c.idCounter, HostDisconnected, serialized)
	c.idCounter++
}

func (c *clientImpl) NotifyTextAdded(id uint64, text string) {
	if c.connection == nil {
		return
	}
	serialized := SerializeTextUpdate(id, text)
	c.connection.SendMessage(c.idCounter, TextUpdate, serialized)
	c.idCounter++
}

func (c *clientImpl) NotifyClientSynced(data *ClientData) {
	if c.connection == nil {
		return
	}
	serialized := SerializeClientData(data)
	c.connection.SendMessage(c.idCounter, HostSynced, serialized)
	c.idCounter++
}

func (c *clientImpl) processFullSyncRequest(id uint64) {
	if c.connection == nil {
		panic("Connection is nil")
	}
	otherClientsData := c.delegate.GetFullSyncData(c)
	serializedd := SerializeSync(c.data, otherClientsData)
	c.connection.SendMessage(id, ServerResponse, serializedd)
}

func (c *clientImpl) processHostSyncRequest(id uint64, data []byte) {
	if c.connection == nil {
		panic("Connection is nil")
	}

	clientId, err := DeserializeClientId(data)
	if err != nil {
		fmt.Printf("Error parsing client ID: %s", err.Error())
		c.reportRequestError(id, "Wrong message sent. Server was unable to parse client ID.")
		return
	}

	clientData := c.delegate.GetClientSyncData(clientId)
	if clientData == nil {
		fmt.Printf("Sync was requested for unknown client %d", clientId)
		c.reportRequestError(id, "Unknown host.")
		return
	}

	serialized := SerializeClientData(clientData)
	c.connection.SendMessage(id, ServerResponse, serialized)
}

func (c *clientImpl) processHostTextUpdate(data []byte) {
	text, err := DeserializeText(data)
	if err != nil {
		fmt.Print("Unable to parse host text update")
		return
	}
	c.data.Data.Text.PushBack(text)
	c.delegate.OnTextAdded(c, text)
}

func (c *clientImpl) processSyncClient(data []byte) {
	clientData, err := DeserializeClientData(data)
	if err != nil {
		fmt.Print("Unable to parse host sync data")
		return
	}
	c.data = clientData
	c.delegate.OnClientSynced(c)
}

func (c *clientImpl) reportRequestError(id uint64, errorText string) {
	if c.connection == nil {
		panic("Connection is nil")
	}
	serialized := SerializeError(errorText)
	c.connection.SendMessage(id, ServerResponse, serialized)
}
