package main

import (
	"bufio"
	cell "container-manager/server/container"
	"container-manager/server/system"
	"fmt"
	"log"
	"net"
)

func main() {

	clients := cell.NewContainerList()
	go system.RunSystem(clients)

	// Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
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

		clients.AddClient(name, conn)
	}
}
