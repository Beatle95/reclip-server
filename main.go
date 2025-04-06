package main

import (
	"errors"
	"fmt"
	"internal"
	"os"
	"slices"
	"strconv"
	"strings"
)

const help = "Server side part of 'Reclip' software. \n" +
	"Arguments: \n" +
	"\t--port=[PORT] (-p [PORT]) - run server on port [PORT] (default value is 8880)\n" +
	"\t--help (-h) - show this help\n"

type parsedSettings struct {
	Port uint16
}

func main() {
	help_pos := slices.IndexFunc(os.Args[1:],
		func(arg string) bool { return arg == "--help" || arg == "-h" })
	if help_pos >= 0 {
		fmt.Print(help)
		return
	}

	settings, err := parseMainArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing command line arguments: %s", err.Error())
	}
	server := internal.CreateServer(settings.Port)
	setClientGroups(server)
	server.Run()
}

func parseMainArgs(args []string) (parsedSettings, error) {
	var settings = parsedSettings{Port: 41286}
	for i := 0; i < len(args); i++ {
		is_port := false
		var port_str string
		if strings.HasPrefix(args[i], "--port=") {
			is_port = true
			port_str = args[i][7:]
		} else if args[i] == "-p" {
			i++
			is_port = true
			if i >= len(args) {
				return settings, errors.New("-p option must be followed by a port number")
			}
			port_str = args[i]
		}

		if is_port {
			value, err := strconv.ParseUint(port_str, 10, 16)
			if err != nil {
				return settings, errors.New("unable to parse port value")
			}
			settings.Port = uint16(value)
		}
	}
	return settings, nil
}

func setClientGroups(server internal.Server) {
	server.AddClientGroup("", internal.CreateClientGroup())
}
