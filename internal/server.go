package internal

import (
	"fmt"
	"log"
	"net"
)

type Server interface {
	Run()
	AddClientGroup(groupId string, newGroup ClientGroup)
}

type serverImpl struct {
	port         uint16
	clientGroups map[string]ClientGroup
}

func CreateServer(port uint16) Server {
	return &serverImpl{
		port: port,
	}
}

func (s *serverImpl) AddClientGroup(groupId string, newGroup ClientGroup) {
	s.clientGroups[groupId] = newGroup
}

func (s *serverImpl) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatal("Error initializing server socket: " + err.Error())
	}

	fmt.Printf("Server up and listening on port %d\n", s.port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *serverImpl) handleConnection(connection net.Conn) {
	log.Printf("Client %v connected.", connection.RemoteAddr())

	// First, client must introduce itself by sending ClientID.
	// TODO:
	// _, err := readClientSecret(&connection)
	// if err != nil {
	// 	connection.Close()
	// 	log.Printf("Unable to read client's ID for client %v, disconnecting...",
	// 		connection.RemoteAddr())
	// 	return
	// }

	err := s.clientGroups[""].AddClient(CreateClient(connection))
	if err != nil {
		fmt.Printf("Unable to add client: %s", connection.RemoteAddr())
		connection.Close()
	}
}
