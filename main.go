package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ProjectOrangeJuice/vm-manager-server/cert"
	"github.com/ProjectOrangeJuice/vm-manager-server/connection"
	"github.com/ProjectOrangeJuice/vm-manager-server/serverConfig"
	"github.com/ProjectOrangeJuice/vm-manager-server/system"
	"github.com/ProjectOrangeJuice/vm-manager-server/web"
)

func main() {
	// if this is the first run, run setup
	config, exists, err := serverConfig.ReadConfig()
	if err != nil {
		if exists {
			log.Printf("Error reading config file, %s. As the file exists, we won't create it", err)
			return
		}
		log.Printf("Config file was not there, running setup [%s]", err)
		err = serverConfig.FirstRun()
		if err != nil {
			log.Printf("First run failed, %s", err)
			return
		}
		return
	}

	log.Printf("Config [%+v]", config)

	tlsConfig, err := cert.SetupTLSConfig("keys/")
	if err != nil {
		log.Fatal(err)
	}

	clientHandler := connection.Setup(config)
	go web.StartWebServer(clientHandler)
	go system.RunSystem(clientHandler)

	listener, err := tls.Listen("tcp", ":8080", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Println("New connection")
		go clientHandler.HandleClient(conn)
	}
}
