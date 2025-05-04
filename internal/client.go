package internal

import "fmt"

type ClientDelegate interface {
	OnClientDisconnected(client Client)
	GetFullSyncData(syncExcluded Client) []ClientData
	GetClientSyncData(id string) *ClientData
	OnTextAdded(client Client, text string)
	OnClientSynced(client Client)
}

type Client interface {
	GetPublicId() string
	GetName() string
	GetClientData() *ClientData
	HandleConnection()

	NotifyClientConnected(id string)
	NotifyClientDisconnected(id string)
	NotifyTextAdded(id string, text string)
	NotifyClientSynced(data *ClientData)
}

type clientImpl struct {
	connection ClientConnection
	delegate   ClientDelegate
	data       ClientData
	idCounter  uint64
}

func CreateClient(
	connection ClientConnection,
	delegate ClientDelegate,
	publicId string,
	name string) Client {
	clientImpl := clientImpl{
		connection: connection,
		delegate:   delegate,
		data: ClientData{
			Id:   publicId,
			Name: name,
		},
		idCounter: 0,
	}
	clientImpl.connection.SetDelegate(&clientImpl)
	return &clientImpl
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
}

// Client implementations:

func (c *clientImpl) GetPublicId() string {
	return c.data.Id
}

func (c *clientImpl) GetName() string {
	return c.data.Name
}

func (c *clientImpl) GetClientData() *ClientData {
	return &c.data
}

func (c *clientImpl) HandleConnection() {
	// Start connection handling.
	c.connection.StartAsync()
}

func (c *clientImpl) NotifyClientConnected(id string) {
	serialized := SerializeClientId(id)
	c.connection.SendMessage(c.idCounter, HostConnected, serialized)
	c.idCounter++
}

func (c *clientImpl) NotifyClientDisconnected(id string) {
	serialized := SerializeClientId(id)
	c.connection.SendMessage(c.idCounter, HostDisconnected, serialized)
	c.idCounter++
}

func (c *clientImpl) NotifyTextAdded(id string, text string) {
	serialized := SerializeTextUpdate(id, text)
	c.connection.SendMessage(c.idCounter, TextUpdate, serialized)
	c.idCounter++
}

func (c *clientImpl) NotifyClientSynced(data *ClientData) {
	serialized := SerializeClientData(data)
	c.connection.SendMessage(c.idCounter, HostSynced, serialized)
	c.idCounter++
}

func (c *clientImpl) processFullSyncRequest(id uint64) {
	otherClientsData := c.delegate.GetFullSyncData(c)
	serializedd := SerializeSync(c.data, otherClientsData)
	c.connection.SendMessage(id, ServerResponse, serializedd)
}

func (c *clientImpl) processHostSyncRequest(id uint64, data []byte) {
	clientId, err := DeserializeClientId(data)
	if err != nil {
		fmt.Printf("Error parsing client ID: %s", err.Error())
		c.reportRequestError(id, "Wrong message sent. Server was unable to parse client ID.")
		return
	}

	clientData := c.delegate.GetClientSyncData(clientId)
	if clientData == nil {
		fmt.Printf("Sync was requested for unknown client %s", clientId)
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
	serialized := SerializeError(errorText)
	c.connection.SendMessage(id, ServerResponse, serialized)
}
