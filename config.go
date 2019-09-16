package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var configfile string = "/etc/stembolt/stembolt.conf"

type Config struct {
	Auth string
	Bind string
}

func ReadConfig() Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file not found: ", configfile)
	}
	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal("error decoding config file: ", err)
	}
	return config
}

func checkAuth(userkey string) bool {
	config := ReadConfig()
	if userkey == config.Auth {
		return true
	} else {
		return false
	}
}
