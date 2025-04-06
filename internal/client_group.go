package internal

import (
	"errors"
)

type ClientGroup interface {
	AddClient(client Client) error
}

type clientGroupImpl struct {
	clients map[string]Client
}

func CreateClientGroup() ClientGroup {
	return &clientGroupImpl{}
}

func (cg *clientGroupImpl) AddClient(client Client) error {
	if cg.clients[client.GetId()] != nil {
		return errors.New("Client with same ID is already registered")
	}
	cg.clients[client.GetId()] = client
	return nil
}
