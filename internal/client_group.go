package internal

import "fmt"

type ClientGroup interface {
	AddClient(client Client)
	RunAsync()
	HandleConnection(id string, connection ClientConnection)

	// ClientDelegate methods:
	GetTaskRunner() EventLoop
	OnClientDisconnected(client Client)
	GetFullSyncData(syncExcluded Client) []ClientData
	GetClientSyncData(id string) *ClientData
	OnTextAdded(client Client, text string)
	OnClientSynced(client Client)
}

type clientGroupImpl struct {
	clients  map[string]Client
	mainLoop EventLoop
	started  bool
}

func CreateClientGroup() ClientGroup {
	return &clientGroupImpl{
		clients:  make(map[string]Client),
		mainLoop: CreateEventLoop(),
		started:  false,
	}
}

// ClientGroup implementations:

func (cg *clientGroupImpl) AddClient(client Client) {
	if cg.started {
		panic("Adding client when group run loop was already started")
	}
	if cg.clients[client.GetClientData().Id] != nil {
		panic("Adding client with existing ID")
	}
	cg.clients[client.GetClientData().Id] = client
}

func (cg *clientGroupImpl) RunAsync() {
	cg.started = true
	go cg.mainLoop.Run()
}

func (cg *clientGroupImpl) HandleConnection(id string, connection ClientConnection) {
	cg.mainLoop.PostTask(
		func() {
			for clientId, client := range cg.clients {
				if clientId == id {
					if !client.IsConnected() {
						client.HandleConnection(connection)
						cg.notifyClientConnected(client.GetClientData().Id)
					} else {
						fmt.Printf("Unable to handle connection from: %s. Client '%s' it already connected.",
							connection.GetAdressString(), client.GetClientData().Name)
						connection.DisconnectAndStop()
					}
					return
				}
			}
			panic("Unable to find corresponding client inside Group")
		},
	)
}

// ClientDelegate implementations:

func (cg *clientGroupImpl) GetTaskRunner() EventLoop {
	return cg.mainLoop
}

func (cg *clientGroupImpl) OnClientDisconnected(client Client) {
	cg.notifyClientDisconnected(client.GetClientData().Id)
}

func (cg *clientGroupImpl) GetFullSyncData(syncExcluded Client) []ClientData {
	var resultData []ClientData
	for clientId, clientValue := range cg.clients {
		if clientId == syncExcluded.GetClientData().Id {
			continue
		}
		resultData = append(resultData, *clientValue.GetClientData())
	}
	return resultData
}

func (cg *clientGroupImpl) GetClientSyncData(id string) *ClientData {
	for clientId, clientValue := range cg.clients {
		if clientId == id {
			return clientValue.GetClientData()
		}
	}
	return nil
}

func (cg *clientGroupImpl) OnTextAdded(client Client, text string) {
	cg.notifyTextAdded(client.GetClientData().Id, text)
}

func (cg *clientGroupImpl) OnClientSynced(client Client) {
	cg.notifyClientSynced(client.GetClientData())
}

func (cg *clientGroupImpl) notifyClientConnected(id string) {
	for clientId, clientValue := range cg.clients {
		if clientId == id {
			continue
		}
		clientValue.NotifyClientConnected(id)
	}
}

func (cg *clientGroupImpl) notifyClientDisconnected(id string) {
	for clientId, clientValue := range cg.clients {
		if clientId == id {
			continue
		}
		clientValue.NotifyClientDisconnected(id)
	}
}

func (cg *clientGroupImpl) notifyTextAdded(id string, text string) {
	for clientId, clientValue := range cg.clients {
		if clientId == id {
			continue
		}
		clientValue.NotifyTextAdded(id, text)
	}
}

func (cg *clientGroupImpl) notifyClientSynced(data *ClientData) {
	for clientId, clientValue := range cg.clients {
		if clientId == data.Id {
			continue
		}
		clientValue.NotifyClientSynced(data)
	}
}
