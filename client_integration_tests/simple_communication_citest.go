package main

import (
	"bufio"
	"communication"
	"fmt"
	"log"
	"os"
)

func RunSimpleCommunicationTest(port uint16) {
	server := communication.CreateServer(port)
	err := server.InitForTesting(3)
	if err != nil {
		fmt.Printf("Unable to initialize the server: '%s'", err.Error())
		return
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
