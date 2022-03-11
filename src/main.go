package main

import (
	"encoding/json"
	"fmt"
	"github.com/melbahja/goph"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Servers struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Site     string `json:"site"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func ExecuteCommand(hostname string, ip string, user string, password string) (error, error) {
	log.Println("Connecting to server: ", hostname)
	client, err := goph.New(user, ip, goph.Password(password))
	if err != nil {
		log.Println("Failed to connect!")
		return nil, err
	}
	out, err := client.Run("pwd && ls -lrt")
	if err != nil {
		log.Println("Failed to execute command!")
		return nil, err
	}
	log.Println(string(out))
	cmd, err := client.Command("ls", "-alh", "/tmp")
	//cmd, err := client.Command("ls", "-alh", "/tmp")
	err = cmd.Run()
	return nil, err
}

func ReadFile(path string) Servers {
	var ArrayServers Servers

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return ArrayServers
	}
	log.Println("Successfully opened ", path)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &ArrayServers)
	if err != nil {
		return Servers{}
	}
	log.Println("Found", len(ArrayServers.Servers), "servers on file")
	//fmt.Println(string(byteValue))
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	return ArrayServers
}

func main() {
	log.Print("Starting GoMonitorSsh...")
	var servers = ReadFile("./resource/server_data.json")
	for i := 0; i < len(servers.Servers); i++ {
		err, _ := ExecuteCommand(servers.Servers[i].Hostname, servers.Servers[i].Ip, servers.Servers[i].User,
			servers.Servers[i].Password)
		if err != nil {
			return
		}
	}
	time.Sleep(15)
}
