package internal

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

const AppName = "reclip-server"
const DefaultServerPort = 41286
const help = "\nServer side part of 'Reclip' software. \n" +
	"Arguments: \n" +
	"\t--port=[PORT] (-p [PORT]) - run server on port [PORT] (default value is 8880)\n" +
	"\t--app-data-dir=[PATH] - override application data directory\n" +
	"\t--help (-h) - show this help\n\n"

type AppSettings struct {
	Port       uint16
	AppDataDir string
}

type ClientConfig struct {
	Secret   string
	PublicId string
	Name     string
}

type GroupConfig struct {
	Clients []ClientConfig
}

type Config struct {
	Groups []GroupConfig
}

func ParseCmdArgs() (AppSettings, error) {
	var port int
	var app_data_dir string

	flag.IntVar(&port, "port", DefaultServerPort, "Run server on port [PORT] (default value is 8880)")
	flag.IntVar(&port, "p", DefaultServerPort, "Run server on port [PORT] (default value is 8880)")
	flag.StringVar(&app_data_dir, "app-data-dir", "", "Override application data directory")
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), help)
	}
	flag.Parse()

	return AppSettings{Port: uint16(port), AppDataDir: app_data_dir}, nil
}

func InitAppDataDir(argsPath string) (string, error) {
	var path string
	if len(argsPath) == 0 {
		defaultPath, _ := os.UserConfigDir()
		path = fmt.Sprintf("%s/%s", defaultPath, AppName)
	} else {
		path = argsPath
	}

	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("application data path is not a directory")
	}

	log.Printf("Server application data directory is: '%s'", path)
	return path, nil
}

func ReadServerConfig(appDataDir string) (*Config, error) {
	configPath := fmt.Sprintf("%s/config.json", appDataDir)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("config is not exists")
	}

	config_data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config")
	}

	var config Config
	err = json.Unmarshal(config_data, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config")
	}
	return &config, nil
}
