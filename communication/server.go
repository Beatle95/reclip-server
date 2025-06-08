package communication

import (
	"errors"
	"fmt"
	"internal"
	"log"
	"net"
)

type secretMapping struct {
	group    internal.ClientGroup
	publicId string
}

type Server interface {
	Init() error
	Run()
}

type serverImpl struct {
	clientGroups  []internal.ClientGroup
	secretMapping map[string]secretMapping

	port        uint16
	initialized bool
}

func CreateServer(port uint16) Server {
	return &serverImpl{
		secretMapping: make(map[string]secretMapping),

		port:        port,
		initialized: false,
	}
}

func (s *serverImpl) Init() error {
	if s.initialized {
		return errors.New("attempting to initialize server twice")
	}

	// TODO: Add groups config and parse it here.
	// Right now we will only support single group.
	new_group := internal.CreateClientGroup()
	for i := 1; i < 5; i++ {
		id := fmt.Sprintf("public%d", i)
		secret := fmt.Sprintf("secret%d", i)
		name := fmt.Sprintf("name%d", i)

		client := internal.CreateClient(new_group, id, name)
		new_group.AddClient(client)
		s.secretMapping[secret] = secretMapping{group: new_group, publicId: id}
	}
	s.clientGroups = append(s.clientGroups, new_group)

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
		new_conn := createClientConnection(conn)
		go s.handleNewConnection(&new_conn)
	}
}

func (s *serverImpl) handleNewConnection(connection *clientConnectionImpl) {
	secretBuf, err := connection.ReadIntroduction()
	if err != nil {
		fmt.Printf("Disconnecting client: %s. Error: %s", connection.GetAdressString(), err.Error())
		return
	}

	secret, err := internal.DeserializeClientIntroduction(secretBuf)
	if err != nil {
		connection.DisconnectAndStop()
		fmt.Printf("Disconnecting client: %s. Error: %s", connection.GetAdressString(), err.Error())
		return
	}

	mapping := s.secretMapping[secret]
	if mapping.group == nil {
		connection.DisconnectAndStop()
		fmt.Printf("Disconnecting unknown client: %s", connection.GetAdressString())
		return
	}

	mapping.group.HandleConnection(mapping.publicId, connection)
}
