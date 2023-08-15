package cell

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

type containerData struct {
	Clients  []Client
	CellLock sync.Mutex
}

type Client struct {
	Name string
	Conn net.Conn
}

type Container interface {
	AddClient(name string, conn net.Conn)
	RemoveClient(name string)
	GetAllClients() []Client
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

	c.Clients = append(c.Clients, Client{
		Name: name,
		Conn: conn,
	})
	go c.processCell(name, conn)
}

func (c *containerData) RemoveClient(name string) {
	c.CellLock.Lock()
	defer c.CellLock.Unlock()
	for i, client := range c.Clients {
		if client.Name == name {
			c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
		}
	}
}

func (c *containerData) GetAllClients() []Client {
	return c.Clients
}

func (c *containerData) processCell(name string, conn net.Conn) {
	// Create a buffered reader
	reader := bufio.NewReader(conn)

	// Read line by line
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
	}
	c.RemoveClient(name)
}
