package main

import (
	"communication"
	"fmt"
	"internal"
	"log"
	"os"
)

const help = "Server side part of 'Reclip' software. \n" +
	"Arguments: \n" +
	"\t--port=[PORT] (-p [PORT]) - run server on port [PORT] (default value is 8880)\n" +
	"\t--app-data-dir=[PATH] - override application data directory" +
	"\t--help (-h) - show this help\n"

func main() {
	if internal.HasHelpArg() {
		fmt.Print(help)
		return
	}

	settings, err := internal.ParseMainArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("Error parsing command line arguments: %s", err.Error())
	}

	appDataDir, err := internal.InitAppDataDir(settings.AppDataDir)
	if err != nil {
		log.Fatalf("Unable to initialize application data directory: '%s'", err.Error())
	}

	config, err := internal.ReadServerConfig(appDataDir)
	if err != nil {
		log.Fatalf("Error parsing server config: '%s'", err.Error())
	}

	server := communication.CreateServer(settings.Port)
	err = server.Init(config)
	if err != nil {
		log.Fatalf("Unable to initialize the server: '%s'", err.Error())
	}
	server.Run()
}
