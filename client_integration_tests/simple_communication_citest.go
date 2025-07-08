package main

import (
	"bufio"
	"communication"
	"log"
	"os"
)

func RunSimpleCommunicationTest(port uint16) {
	server, err := communication.CreateServerForTesting(port, 3)
	if err != nil {
		log.Fatalf("Unable to initialize the server: '%s'", err.Error())
	}
	go waitStop()
	server.Run()
}

func waitStop() {
	reader := bufio.NewReader(os.Stdin)
	command, _ := reader.ReadString('\n')
	if command == "stop\n" {
		os.Exit(0)
	} else {
		log.Fatalf("Unknown command received: %s", command)
	}
}
