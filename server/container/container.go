package cell

import (
	"bufio"
	"container-manager/shared"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type containerData struct {
	Clients  []*Client
	CellLock sync.Mutex
}

type Client struct {
	Name string
	Conn net.Conn

	// The client information
	Storage []shared.StorageResult
	System  shared.SystemResult
}

type Container interface {
	AddClient(name string, conn net.Conn)
	RemoveClient(name string)
	GetAllClients() []*Client
}

func NewContainerList() Container {
	return &containerData{}
}

func (c *containerData) AddClient(name string, conn net.Conn) {
	c.CellLock.Lock()
	defer c.CellLock.Unlock()
	// Check if the name already exists
	for _, client := range c.Clients {
		if client.Name == name {
			fmt.Fprint(conn, "Name already exists\n")
			conn.Close()
			return
		}
	}
	newClient := Client{
		Name: name,
		Conn: conn,
	}
	log.Printf("Adding client %s", name)
	c.Clients = append(c.Clients, &newClient)
	go c.processCell(&newClient)

	fmt.Fprint(newClient.Conn, "STORAGE_INFO\n")
	fmt.Fprint(newClient.Conn, "SYSTEM_INFO\n")
}

func (c *containerData) RemoveClient(name string) {
	c.CellLock.Lock()
	log.Printf("Removing client %s", name)
	defer c.CellLock.Unlock()
	for i, client := range c.Clients {
		if client.Name == name {
			c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
		}
	}
}

func (c *containerData) GetAllClients() []*Client {
	return c.Clients
}

func (c *containerData) processCell(client *Client) {
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
	c.RemoveClient(client.Name)
}
