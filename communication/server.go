package communication

import (
	"crypto/tls"
	"fmt"
	"internal"
	"log"
	"os"
	"path/filepath"
)

type secretMapping struct {
	group    internal.ClientGroup
	publicId uint64
}

type Server struct {
	tlsConfig     *tls.Config
	clientGroups  []internal.ClientGroup
	secretMapping map[string]secretMapping
	port          uint16
}

func CreateServer(appDataDir string, port uint16, appConfig *internal.Config) (*Server, error) {
	result := &Server{
		secretMapping: make(map[string]secretMapping),
		port:          port,
	}

	var err error
	result.tlsConfig, err = loadTlsConfig(appDataDir)
	if err != nil {
		return nil, fmt.Errorf("unable to load TLS config: %v", err)
	}

	for _, groupConfig := range appConfig.Groups {
		newGroup := internal.CreateClientGroup()
		for _, clientConfig := range groupConfig.Clients {
			client := internal.CreateClient(newGroup, clientConfig.PublicId, clientConfig.Name)
			if _, exists := result.secretMapping[clientConfig.Secret]; exists {
				return nil, fmt.Errorf("initialization error, there was multiple clients " +
					"with the same secret ID in the config")
			}
			result.secretMapping[clientConfig.Secret] = secretMapping{
				group: newGroup, publicId: clientConfig.PublicId}
			newGroup.AddClient(client)
			result.clientGroups = append(result.clientGroups, newGroup)
		}
	}

	return result, nil
}

func CreateServerForTesting(port uint16, clients_count int) (*Server, error) {
	result := &Server{
		secretMapping: make(map[string]secretMapping),
		port:          port,
	}

	var err error
	result.tlsConfig, err = LoadTestTlsConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load TLS config: %v", err)
	}

	new_group := internal.CreateClientGroup()
	for i := 1; i <= clients_count; i++ {
		id := uint64(i)
		secret := fmt.Sprintf("secret%d", i)
		name := fmt.Sprintf("name%d", i)

		client := internal.CreateClient(new_group, id, name)
		new_group.AddClient(client)
		result.secretMapping[secret] = secretMapping{group: new_group, publicId: id}
	}
	result.clientGroups = append(result.clientGroups, new_group)

	return result, nil
}

func (s *Server) Run() {
	for _, group := range s.clientGroups {
		group.RunAsync()
	}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", s.port), s.tlsConfig)
	if err != nil {
		log.Fatal("Error initializing server socket: " + err.Error())
	}

	defer listener.Close()
	log.Printf("Server up and listening on port %d\n", s.port)
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

func (s *Server) handleNewConnection(connection *clientConnectionImpl) {
	secretBuf, err := connection.ReadIntroduction()
	if err != nil {
		log.Printf("Disconnecting client: %s. Error: %s", connection.GetAdressString(), err.Error())
		return
	}

	secret, err := internal.DeserializeClientIntroduction(secretBuf)
	if err != nil {
		connection.DisconnectAndStop()
		log.Printf("Disconnecting client: %s. Error: %s", connection.GetAdressString(), err.Error())
		return
	}

	mapping, mappingExists := s.secretMapping[secret]
	if !mappingExists {
		connection.DisconnectAndStop()
		log.Printf("Disconnecting unknown client: %s", connection.GetAdressString())
		return
	}

	mapping.group.HandleConnection(mapping.publicId, connection)
}

func loadTlsConfig(certDir string) (*tls.Config, error) {
	certFile := certDir + "/cert.pem"
	keyFile := certDir + "/key.pem"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return certToConfig(&cert), nil
}

func LoadTestTlsConfig() (*tls.Config, error) {
	bin_path, err := os.Executable()
	if err != nil {
		return nil, err
	}

	test_data_dir := filepath.Dir(bin_path) + "/../resources/test_data"
	certFile := test_data_dir + "/cert.pem"
	keyFile := test_data_dir + "/key.pem"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return certToConfig(&cert), nil
}

func certToConfig(cert *tls.Certificate) *tls.Config {
	config := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		},
	}
	return config
}
