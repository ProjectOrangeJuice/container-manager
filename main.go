package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ProjectOrangeJuice/vm-manager-server/cert"
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

		// Get the client's name
		nameReader := bufio.NewReader(conn)
		name, err := nameReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		log.Printf("line -> %s", name)
		// clients.AddClient(name, conn)
	}
}
