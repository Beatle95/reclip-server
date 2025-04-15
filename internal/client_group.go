package internal

import (
	"errors"
)

type ClientGroup interface {
	AddClient(client Client) error
	OnClientDisconnected(id string)
}

type clientGroupImpl struct {
	clients map[string]Client
}

func CreateClientGroup() ClientGroup {
	return &clientGroupImpl{}
}

func (cg *clientGroupImpl) AddClient(client Client) error {
	if cg.clients[client.GetPublicId()] != nil {
		return errors.New("Client with same ID is already registered")
	}
	cg.clients[client.GetPublicId()] = client
	cg.notifyClientConnected(client.GetPublicId())
	return nil
}

func (cg *clientGroupImpl) OnClientDisconnected(id string) {
	// TODO:
}

func (cg *clientGroupImpl) notifyClientConnected(id string) {
	// TODO:
}
