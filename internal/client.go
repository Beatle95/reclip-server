package internal

type ClientDelegate interface {
	OnClientDisconnected(client Client)
	// OnTextAdded(id string, text string)
	// OnClientSynced()
}

type Client interface {
	GetPublicId() string
	GetName() string
	HandleConnection()
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
	clientImpl := clientImpl{
		connection: connection,
		delegate:   delegate,
		data: ClientData{
			id:   publicId,
			name: name,
		},
	}
	clientImpl.connection.SetDelegate(&clientImpl)
	return &clientImpl
}

// ClientConnectionDelegate implementations:

func (c *clientImpl) ProcessMessage(id uint64, msgType ClientMessageType, data []byte) {
	// TODO:
}

func (c *clientImpl) OnDisconnected() {
	c.delegate.OnClientDisconnected(c)
}

// Client implementations:

func (c *clientImpl) GetPublicId() string {
	return c.data.id
}

func (c *clientImpl) GetName() string {
	return c.data.name
}

func (c *clientImpl) HandleConnection() {
	// Start connection handling.
	c.connection.StartAsync()
}
