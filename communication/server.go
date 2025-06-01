package communication

import (
	"errors"
	"fmt"
	"internal"
	"log"
	"net"
)

type Server interface {
	Init() error
	Run()
}

type serverImpl struct {
	port         uint16
	clientGroups map[string]internal.ClientGroup
	initialized  bool
}

func CreateServer(port uint16) Server {
	return &serverImpl{
		port:         port,
		clientGroups: make(map[string]internal.ClientGroup),
		initialized:  false,
	}
}

func (s *serverImpl) Init() error {
	if s.initialized {
		return errors.New("attempting to initialize server twice")
	}
	s.clientGroups[""] = internal.CreateClientGroup()

	for _, group := range s.clientGroups {
		group.RunAsync()
	}

	s.initialized = true
	return nil
}

func (s *serverImpl) Run() {
	if !s.initialized {
		log.Fatal("Server is not initialized")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatal("Error initializing server socket: " + err.Error())
	}

	defer listener.Close()
	fmt.Printf("Server up and listening on port %d\n", s.port)
	for {
		conn, err := listener.Accept()
		log.Printf("Client %v connected.", conn.RemoteAddr())
		if err != nil {
			log.Println(err)
			continue
		}
		HandleUnregisteredClient(conn, s.handleNewConnection)
	}
}

// Called from another goroutine.
func (s *serverImpl) handleNewConnection(connection net.Conn, secret string) {
	group := s.findGroupForClient(secret)
	if group == nil {
		connection.Close()
		fmt.Printf("Disconnecting unknown client: %s", connection.RemoteAddr().String())
		return
	}
	connInternal := createClientConnection(connection)
	client := internal.CreateClient(&connInternal, group, "public", "name")
	err := group.HandleNewClient(client)
	if err != nil {
		fmt.Printf("Unable to add client: %s", connInternal.GetAdressString())
		connection.Close()
	}
}

func (s *serverImpl) findGroupForClient(clientSecret string) internal.ClientGroup {
	return s.clientGroups[clientSecret]
}
