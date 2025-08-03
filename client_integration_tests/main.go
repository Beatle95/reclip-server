package main

import (
	"communication"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const TimeoutSec = 45

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("Usage: client_integration_tests <test_name> <port>")
	}

	test_name := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("Unable parse port: %s", err.Error())
	}
	if port < 0 || port > math.MaxUint16 {
		log.Fatalf("Unexeptable port value: %s", args[1])
	}

	go startTimeoutWatchdog()
	runTest(test_name, uint16(port))
}

func runTest(test_name string, port uint16) {
	switch test_name {
	case "communication_protocol_test":
		communication.RunCommunicationProtocolTest(port)
	case "simple_communication_test":
		RunSimpleCommunicationTest(port, 3)
	case "multiple_communication_test":
		RunSimpleCommunicationTest(port, 30)
	default:
		log.Fatalf("Unknown test: %s", test_name)
	}
}

func startTimeoutWatchdog() {
	start := time.Now()
	for time.Since(start) < time.Second*TimeoutSec {
		time.Sleep(time.Millisecond * 100)
	}

	if time.Since(start) > time.Second*TimeoutSec {
		panic("Test timeout")
	}
}
