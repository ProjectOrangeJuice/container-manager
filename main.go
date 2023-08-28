package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ProjectOrangeJuice/vm-manager-server/cert"
	"github.com/ProjectOrangeJuice/vm-manager-server/connection"
)

func main() {

	// clients := cell.NewContainerList()
	// go system.RunSystem(clients)
	// go web.StartWebServer(clients)
	// Listen on port 8080

	tlsConfig, err := cert.SetupTLSConfig("keys/")
	if err != nil {
		log.Fatal(err)
	}

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
		go connection.HandleClient(conn)
		// clients.AddClient(name, conn)
	}
}
