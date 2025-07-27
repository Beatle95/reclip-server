package main

import (
	"communication"
	"internal"
	"log"
)

func main() {
	settings, err := internal.ParseCmdArgs()
	if err != nil {
		log.Fatalf("Error parsing command line arguments: %s", err.Error())
	}

	appDataDir, err := internal.InitAppDataDir(settings.AppDataDir)
	if err != nil {
		log.Fatalf("Unable to initialize application data directory: '%s'", err.Error())
	}
	log.Printf("Server application data directory is: '%s'", appDataDir)

	config, err := internal.ReadServerConfig(appDataDir)
	if err != nil {
		log.Fatalf("Error parsing server config: '%s'", err.Error())
	}

	server, err := communication.CreateServer(settings.AppDataDir, settings.Port, config)
	if err != nil {
		log.Fatalf("Unable to initialize the server: '%s'", err.Error())
	}
	server.Run()
}
