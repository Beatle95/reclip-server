package internal

import (
	"errors"
	"sync"
)

type ClientGroup interface {
	HandleNewClient(client Client) error

	// ClientDelegate methods:
	OnClientDisconnected(client Client)
	GetFullSyncData(syncExcluded Client) []ClientData
	GetClientSyncData(id string) *ClientData
	OnTextAdded(client Client, text string)
	OnClientSynced(client Client)
}

type clientGroupImpl struct {
	clientsMapMutex sync.Mutex
	clients         map[string]Client
}

func CreateClientGroup() ClientGroup {
	return &clientGroupImpl{
		clients: make(map[string]Client),
	}
}

// ClientGroup implementations:

func (cg *clientGroupImpl) HandleNewClient(client Client) error {
	cg.clientsMapMutex.Lock()
	defer cg.clientsMapMutex.Unlock()

	if cg.clients[client.GetPublicId()] != nil {
		return errors.New("Client with same ID is already registered")
	}
	cg.clients[client.GetPublicId()] = client
	client.HandleConnection()
	cg.notifyClientConnected(client.GetPublicId())
	return nil
}

// ClientDelegate implementations:

func (cg *clientGroupImpl) OnClientDisconnected(client Client) {
	cg.clientsMapMutex.Lock()
	defer cg.clientsMapMutex.Unlock()

	cg.clients[client.GetPublicId()] = nil
	cg.notifyClientDisconnected(client.GetPublicId())
}

func (cg *clientGroupImpl) GetFullSyncData(syncExcluded Client) []ClientData {
	cg.clientsMapMutex.Lock()
	defer cg.clientsMapMutex.Unlock()

	var resultData []ClientData
	for clientId, clientValue := range cg.clients {
		if clientId == syncExcluded.GetPublicId() {
			continue
		}
		resultData = append(resultData, *clientValue.GetClientData())
	}
	return resultData
}

func (cg *clientGroupImpl) GetClientSyncData(id string) *ClientData {
	cg.clientsMapMutex.Lock()
	defer cg.clientsMapMutex.Unlock()

	for clientId, clientValue := range cg.clients {
		if clientId == id {
			return clientValue.GetClientData()
		}
	}
	return nil
}

func (cg *clientGroupImpl) OnTextAdded(client Client, text string) {
	cg.clientsMapMutex.Lock()
	defer cg.clientsMapMutex.Unlock()
	cg.notifyTextAdded(client.GetPublicId(), text)
}

func (cg *clientGroupImpl) OnClientSynced(client Client) {
	cg.clientsMapMutex.Lock()
	defer cg.clientsMapMutex.Unlock()
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
