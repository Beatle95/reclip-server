package internal

import "net"

type Client interface {
	GetId() string
}

func CreateClient(connection net.Conn) Client {
	return nil
}

// type ClientImpl struct {
// 	Connection net.Conn
// 	ClientId   string

// 	group            *ClientGroup
// 	name             string
// 	clipboardContent []string
// }

// func CreateClient(group *ClientGroup, connection net.Conn, client_id string) *Client {
// 	return &Client{
// 		Connection: connection,
// 		ClientId:   client_id,
// 		group:      group,
// 	}
// }

// func Handle(self *Client) {
// 	// TODO:
// 	for {
// 	}
// }
