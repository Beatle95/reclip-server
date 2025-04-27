package communication

import "net"

type RegistrationCallback func(connection net.Conn, secret string)

type unknownClientConnection struct {
	connection net.Conn
	callback   RegistrationCallback
}

func HandleUnregisteredClient(conn net.Conn, callback RegistrationCallback) {
	client := unknownClientConnection{
		connection: conn,
		callback:   callback,
	}
	go client.run()
}

func (c *unknownClientConnection) run() {
	// TODO:
	c.callback(c.connection, "")
}
