package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const version string = "v0.0.1"

type Config struct {
	ServerPort     int    `json:"server-port"`
	StoreDirectory string `json:"store-dir"`
}

func main() {
	conf := readConf()
	store, err := newStore(conf.StoreDirectory)
	if err != nil {
		log.Fatalln(err)
	}
	r := NewRoutes(store)

	url := fmt.Sprintf("localhost:%d", conf.ServerPort)
	log.Printf("Starting server on port %d..\n", conf.ServerPort)
	r.Run(url)
}

func readConf() *Config {
	// open json file
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	// defer the closing of file
	defer jsonFile.Close()

	// read file as byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// initialize & unmarshal config
	var conf Config
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		log.Fatalln("Error reading config")
	}

	return &conf
}
