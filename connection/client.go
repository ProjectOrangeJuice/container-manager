package connection

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/ProjectOrangeJuice/vm-manager-server/shared"
)

type allClients struct {
	Clients    []*Client
	clientLock sync.Mutex
}

type Client struct {
	Name   string
	Serial string
	Conn   net.Conn

	// The client information
	Storage []shared.StorageResult
	System  shared.SystemResult
}

// AddClient adds a new client to the list of clients and reads the data from the client connection
func (ac *allClients) AddClient(name, serial string, conn net.Conn) {
	ac.clientLock.Lock()
	defer ac.clientLock.Unlock()
	// Check if the name already exists
	for _, client := range ac.Clients {
		if client.Name == name {
			fmt.Fprint(conn, "Name already exists\n")
			conn.Close()
			return
		}
	}
	newClient := Client{
		Name:   name,
		Conn:   conn,
		Serial: serial,
	}
	log.Printf("Adding client %s", name)
	ac.Clients = append(ac.Clients, &newClient)
	go ac.processCell(&newClient)

	// Send the client the information it needs
	fmt.Fprint(newClient.Conn, "STORAGE_INFO\n")
	fmt.Fprint(newClient.Conn, "SYSTEM_INFO\n")
}

func (ac *allClients) RemoveClient(name string) {
	ac.clientLock.Lock()
	log.Printf("Removing client %s", name)
	defer ac.clientLock.Unlock()
	for i, client := range ac.Clients {
		if client.Name == name {
			ac.Clients = append(ac.Clients[:i], ac.Clients[i+1:]...)
		}
	}
}

func (ac *allClients) GetAllClients() []*Client {
	return ac.Clients
}

func (ac *allClients) processCell(client *Client) {
	// Create a buffered reader
	reader := bufio.NewReader(client.Conn)

	for {
		// Read a line of data
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		// Print the line
		fmt.Println(line)
		client.processEvent(line)
	}
	ac.RemoveClient(client.Name)
}
