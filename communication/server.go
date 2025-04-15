package communication

import (
	"fmt"
	"internal"
	"log"
	"net"
)

type Server interface {
	Run()
	AddClientGroup(groupId string, newGroup internal.ClientGroup)
}

type serverImpl struct {
	port         uint16
	clientGroups map[string]internal.ClientGroup
}

func CreateServer(port uint16) Server {
	return &serverImpl{
		port: port,
	}
}

func (s *serverImpl) AddClientGroup(groupId string, newGroup internal.ClientGroup) {
	s.clientGroups[groupId] = newGroup
}

func (s *serverImpl) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatal("Error initializing server socket: " + err.Error())
		return
	}

	defer listener.Close()
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

	// First, client must introduce itself by sending secret ID.
	clientSecret, error := s.processIntroduction(connection)
	if error != nil {
		log.Printf("Client %v introduction error: %v", connection.RemoteAddr(), error)
		connection.Close()
		return
	}

	group := s.clientGroups[clientSecret]
	client := internal.CreateClient(
		&clientConnectionImpl{connection: connection},
		group, "public", "name")
	err := group.AddClient(client)
	if err != nil {
		fmt.Printf("Unable to add client: %s", connection.RemoteAddr())
		connection.Close()
	}

	client.HandleConnection()
}

func (s *serverImpl) processIntroduction(_ net.Conn) (string, error) {
	// TODO:
	return "secret", nil
}
