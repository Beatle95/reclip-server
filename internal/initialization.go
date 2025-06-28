package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const AppName = "reclip-server"
const DefaultServerPort = 41286

type AppSettings struct {
	Port       uint16
	AppDataDir string
}

type GroupConfig struct {
	Clients []ClientConfig
}

type ClientConfig struct {
	Secret   string
	PublicId string
	Name     string
}

type Config struct {
	Groups []GroupConfig
}

func HasHelpArg() bool {
	help_pos := slices.IndexFunc(os.Args[1:],
		func(arg string) bool { return arg == "--help" || arg == "-h" })
	return help_pos > 0
}

func ParseMainArgs(args []string) (AppSettings, error) {
	var settings = AppSettings{Port: DefaultServerPort}
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
		} else if strings.HasPrefix(args[i], "--app-data-dir=") {
			settings.AppDataDir = args[i][15:]
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

func InitAppDataDir(arsgPath string) (string, error) {
	var path string
	if len(arsgPath) == 0 {
		defaultPath, _ := os.UserConfigDir()
		path = fmt.Sprintf("%s/%s", defaultPath, AppName)
	} else {
		path = arsgPath
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
