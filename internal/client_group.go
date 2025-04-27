package internal

import (
	"errors"
	"sync"
)

type ClientGroup interface {
	HandleNewClient(client Client) error
	OnClientDisconnected(client Client)
}

type clientGroupImpl struct {
	clientsMutex sync.Mutex
	clients      map[string]Client
}

func CreateClientGroup() ClientGroup {
	return &clientGroupImpl{
		clients: make(map[string]Client),
	}
}

// ClientGroup implementations:

func (cg *clientGroupImpl) HandleNewClient(client Client) error {
	cg.clientsMutex.Lock()
	defer cg.clientsMutex.Unlock()

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
	cg.clientsMutex.Lock()
	defer cg.clientsMutex.Unlock()

	cg.clients[client.GetPublicId()] = nil
	cg.notifyClientDisconnected(client.GetPublicId())
}

func (cg *clientGroupImpl) notifyClientConnected(id string) {
	// TODO:
}

func (cg *clientGroupImpl) notifyClientDisconnected(id string) {
	// TODO:
}
