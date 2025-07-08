package communication

import (
	"crypto/tls"
	"fmt"
	"internal"
	"log"
)

type TestConnectionDelegate struct {
	connection    internal.ClientConnection
	eventLoop     internal.EventLoop
	clientMsgType internal.ClientMessageType
	serverMsgType internal.ServerMessageType
}

func RunCommunicationProtocolTest(port uint16) {
	config, err := LoadTestTlsConfig()
	if err != nil {
		log.Fatalf("Unable to initialize test tls config: %s", err.Error())
	}

	log.Printf("Starting listening on port %d", port)
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)
	if err != nil {
		log.Fatal("Error initializing server socket: " + err.Error())
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Client %v connected.", conn.RemoteAddr())
	defer conn.Close()

	clientConnection := CreateClientConnectionForTesting(conn)
	testDelegate := TestConnectionDelegate{
		connection:    clientConnection,
		eventLoop:     internal.CreateEventLoop(),
		clientMsgType: internal.ClientResponse,
		serverMsgType: internal.ServerResponse,
	}

	clientConnection.SetUp(&testDelegate, testDelegate.eventLoop)
	clientConnection.StartHandlingAsync()
	testDelegate.eventLoop.Run()
	clientConnection.DisconnectAndStop()
}

func (test *TestConnectionDelegate) OnDisconnected() {
	test.eventLoop.Quit()
}

func (test *TestConnectionDelegate) ProcessMessage(
	id uint64, msgType internal.ClientMessageType, data []byte) {
	if msgType != test.clientMsgType {
		log.Fatalf("Received unexpected client message type: %d", int(msgType))
	}
	test.connection.SendMessage(id, test.serverMsgType, data)
	test.clientMsgType = incrementClientMessageType(test.clientMsgType)
	test.serverMsgType = incrementServerMessageType(test.serverMsgType)
}

func incrementClientMessageType(val internal.ClientMessageType) internal.ClientMessageType {
	val += 1
	if val > internal.ClientMessageTypeMax {
		return internal.ClientResponse
	} else {
		return val
	}
}

func incrementServerMessageType(val internal.ServerMessageType) internal.ServerMessageType {
	val += 1
	if val > internal.ServerMessageTypeMax {
		return internal.ServerResponse
	} else {
		return val
	}
}
